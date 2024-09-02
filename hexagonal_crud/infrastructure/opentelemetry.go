package infrastructure

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"

	"go.opentelemetry.io/otel"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	otelTrace "go.opentelemetry.io/otel/trace"
)

// SetupOTelSDK mengatur OpenTelemetry SDK untuk aplikasi dengan menginisialisasi tracer dan meter provider.
// Fungsi ini mengembalikan fungsi shutdown untuk membersihkan sumber daya saat aplikasi berhenti.
func SetupOTelSDK(ctx context.Context, endpoint, serviceName, serviceVersion, environment string) (shutdown func(context.Context) error, err error) {
	// shutdownFuncs menyimpan fungsi-fungsi yang perlu dijalankan saat aplikasi berhenti.
	var shutdownFuncs []func(context.Context) error

	// Fungsi shutdown yang akan memanggil semua fungsi shutdown yang disimpan dalam shutdownFuncs.
	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	// Fungsi untuk menangani kesalahan dan melakukan shutdown jika terjadi error.
	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	// Mengatur resource yang digunakan oleh OpenTelemetry.
	res, err := newResource(serviceName, serviceVersion, environment)
	if err != nil {
		handleErr(err)
		return
	}

	// Mengatur propagator untuk context propagation di OpenTelemetry.
	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	// Mengatur tracer provider untuk tracing.
	tracerProvider, err := newTraceProvider(ctx, endpoint, res)
	if err != nil {
		handleErr(err)
		return
	}
	// Menambahkan fungsi shutdown untuk tracer provider.
	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)

	// Mengatur meter provider untuk metrik.
	meterProvider, err := newMeterProvider(ctx, endpoint, res)
	if err != nil {
		handleErr(err)
		return
	}
	// Menambahkan fungsi shutdown untuk meter provider.
	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)

	return
}

// newResource membuat resource baru yang berisi informasi tentang layanan seperti nama, versi, dan environment.
func newResource(serviceName, serviceVersion, environment string) (*resource.Resource, error) {
	return resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(serviceVersion),
			semconv.DeploymentEnvironment(environment),
		))
}

// newPropagator mengatur propagator yang digunakan untuk context propagation dalam OpenTelemetry.
func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

// newTraceProvider mengatur trace provider menggunakan OTLP (OpenTelemetry Protocol) untuk mengirim data trace.
func newTraceProvider(ctx context.Context, endpoint string, res *resource.Resource) (*trace.TracerProvider, error) {

	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure(), otlptracegrpc.WithEndpoint(endpoint))
	if err != nil {
		log.Fatalf("new otlp trace grpc exporter failed: %v", err)
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter,
			// Mengatur batch timeout agar setiap batch dikirim setiap 1 detik.
			trace.WithBatchTimeout(time.Second)),
		trace.WithResource(res),
	)
	return traceProvider, nil
}

// newMeterProvider mengatur meter provider menggunakan OTLP untuk mengirim data metrik.
func newMeterProvider(ctx context.Context, endpoint string, res *resource.Resource) (*metric.MeterProvider, error) {

	exporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithInsecure(), otlpmetricgrpc.WithEndpoint(endpoint))
	if err != nil {
		log.Fatalf("new otlp metric grpc exporter failed: %v", err)
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(exporter,
			// Mengatur interval pengiriman metrik setiap 3 detik.
			metric.WithInterval(3*time.Second))),
	)
	return meterProvider, nil
}

// Tracer mengembalikan Tracer yang diinisialisasi dengan nama layanan yang diambil dari environment variable.
func Tracer() otelTrace.Tracer {
	return otel.Tracer(os.Getenv("OTEL_SERVICE_NAME"))
}

/*
Penjelasan Fungsi dan Komentar dalam Kode:
Package infrastructure:

Package ini bertujuan untuk mengatur infrastruktur OpenTelemetry dalam aplikasi, seperti konfigurasi tracing dan metrik.
Fungsi SetupOTelSDK:

Fungsi ini menginisialisasi OpenTelemetry SDK, termasuk tracer dan meter provider, berdasarkan endpoint dan informasi layanan yang diberikan.
shutdownFuncs menyimpan fungsi-fungsi shutdown yang akan dijalankan ketika aplikasi berhenti.
Mengatur resource untuk tracer dan meter provider yang mengandung metadata seperti nama layanan, versi, dan environment.
propagator digunakan untuk context propagation, yang memastikan trace context dapat diteruskan antar layanan.
Tracer provider dan meter provider diatur dan disimpan, sehingga dapat digunakan di seluruh aplikasi.
Fungsi newResource:

Fungsi ini membuat resource baru yang menyimpan informasi seperti nama layanan, versi, dan environment. Resource ini digunakan oleh tracer dan meter provider.
Fungsi newPropagator:

Fungsi ini mengembalikan propagator yang digunakan untuk meneruskan konteks tracing antar komponen dalam sistem.
Fungsi newTraceProvider:

Fungsi ini membuat trace provider yang digunakan untuk mengirim data tracing ke endpoint OTLP. Data dikirim dalam batch, dengan timeout batch diatur untuk mengirim data setiap 1 detik.
Fungsi newMeterProvider:

Fungsi ini membuat meter provider yang digunakan untuk mengirim data metrik ke endpoint OTLP. Data metrik dikirim setiap 3 detik.
Fungsi Tracer:

Fungsi ini mengembalikan Tracer yang diinisialisasi dengan nama layanan yang diambil dari environment variable. Tracer digunakan untuk membuat span baru dalam trace.
Komentar yang diberikan bertujuan untuk membantu pemahaman tentang bagaimana fungsi-fungsi ini bekerja, dan bagaimana mereka diintegrasikan untuk menyediakan observabilitas dalam aplikasi.
*/
