package main

import (
	"CRUD_Hexagonal/api/product"                  // Mengimpor handler untuk produk
	_ "CRUD_Hexagonal/api/product"                // Mengimpor handler (dengan underscore untuk side effect)
	"CRUD_Hexagonal/infrastructure"               // Mengimpor setup infrastruktur
	storeRepo "CRUD_Hexagonal/repository/product" // Mengimpor repository produk
	storeServ "CRUD_Hexagonal/service/product"    // Mengimpor service produk
	"context"                                     // Mengimpor context untuk manajemen lifecycle aplikasi
	"errors"                                      // Mengimpor errors untuk manajemen error

	otelfiber "github.com/gofiber/contrib/otelfiber/v2"              // Mengimpor middleware OpenTelemetry untuk Fiber
	"github.com/gofiber/fiber/v2"                                    // Mengimpor framework Fiber
	middlewareLogger "github.com/gofiber/fiber/v2/middleware/logger" // Mengimpor middleware logger untuk Fiber
	_ "github.com/joho/godotenv/autoload"                            // Mengimpor dan mengautoload variabel lingkungan dari file .env
	"github.com/spf13/viper"                                         // Mengimpor Viper untuk manajemen konfigurasi

	"log" // Mengimpor log untuk logging standar
	"os"  // Mengimpor os untuk akses sistem operasi, seperti variabel lingkungan

	"golang.org/x/exp/slog" // Mengimpor slog untuk logging yang efisien
)

func main() {
	// Menginisialisasi Viper untuk membaca variabel lingkungan
	viper.AutomaticEnv()
	ctx := context.Background() // Membuat context dasar untuk aplikasi

	// Konfigurasi logger untuk aplikasi dengan JSON output
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo, // Set level log ke info
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			return a // Tidak ada atribut yang diganti
		},
	}).WithAttrs([]slog.Attr{
		slog.String("service", os.Getenv("OTEL_SERVICE_NAME")), // Menambahkan nama layanan dari env var
		slog.String("with-release", "v1.0.0"),                  // Menambahkan versi layanan
	})
	logger := slog.New(logHandler) // Membuat logger baru
	slog.SetDefault(logger)        // Set logger default

	// Mengatur OpenTelemetry untuk observabilitas
	serviceName := os.Getenv("OTEL_SERVICE_NAME") // Mendapatkan nama layanan dari variabel lingkungan
	otelCollector := os.Getenv("OTEL_COLLECTOR")  // Mendapatkan alamat kolektor OTel dari variabel lingkungan
	serviceVersion := "0.1.0"                     // Versi layanan
	otelShutdown, err := infrastructure.SetupOTelSDK(ctx, otelCollector, serviceName, serviceVersion, os.Getenv("OTEL_ENV"))
	if err != nil {
		log.Fatalf("failed to initialize OTel SDK: %v", err) // Jika setup OTel gagal, hentikan aplikasi
	}
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background())) // Menutup OTel dengan benar
	}()

	// Inisialisasi MongoDB
	mongo := infrastructure.NewMongo(ctx, os.Getenv("MONGO_DSN"), os.Getenv("MONGO_DB_NAME"))
	mongo = mongo.Connect() // Menghubungkan ke MongoDB

	// Inisialisasi repository dan service untuk produk
	storeRepository := storeRepo.NewstoreRepository(mongo.Client, mongo.DB, "products") // Membuat repository untuk produk
	storeService := storeServ.NewStoreService(storeRepository)                          // Membuat service produk dengan menggunakan repository
	handler := product.NewStoreHandler(storeService)                                    // Membuat handler untuk produk dengan menggunakan service

	// Inisialisasi aplikasi Fiber
	app := fiber.New()

	// Middleware untuk logging request
	app.Use(middlewareLogger.New(middlewareLogger.Config{
		Format: "[${time}] ${ip}  ${status} - ${latency} ${method} ${path}\n", // Format log request
	}))

	// Middleware untuk integrasi OpenTelemetry dengan Fiber
	app.Use(otelfiber.Middleware())

	// Route untuk homepage
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!") // Mengirim string "Hello, World!" sebagai response
	})

	// Route untuk CRUD API produk
	app.Get("/product/:id", handler.Get)       // Mendapatkan produk berdasarkan ID
	app.Get("/product", handler.GetAll)        // Mendapatkan semua produk
	app.Post("/product", handler.Create)       // Membuat produk baru
	app.Put("/product/:id", handler.Update)    // Memperbarui produk berdasarkan ID
	app.Delete("/product/:id", handler.Delete) // Menghapus produk berdasarkan ID

	// Menjalankan server pada port 3000
	app.Listen(":3000")
}

/*
Fungsi dan Kode yang Dijelaskan:
Viper Initialization (viper.AutomaticEnv):

Viper digunakan untuk membaca variabel lingkungan dari sistem. Ini berguna untuk konfigurasi aplikasi tanpa harus hardcode nilai-nilai tertentu.
Logger Configuration:

Bagian ini mengatur logger menggunakan slog, dengan output JSON. Ini berguna untuk melacak aktivitas aplikasi, terutama dalam lingkungan produksi.
OpenTelemetry Setup:

OpenTelemetry diatur untuk membantu dalam observabilitas aplikasi (melacak performa, logging, tracing). Jika ada kesalahan dalam pengaturan, aplikasi akan berhenti.
MongoDB Initialization:

Bagian ini menginisialisasi koneksi ke MongoDB menggunakan detail yang diberikan melalui variabel lingkungan. Koneksi ini kemudian digunakan untuk mengakses database.
Repository and Service Initialization:

Membuat storeRepository untuk menangani akses data produk dan storeService untuk mengelola logika bisnis terkait produk. handler kemudian digunakan untuk menghubungkan service dengan API.
Fiber Web Framework Setup:

Fiber adalah framework web yang digunakan untuk membangun API ini. Middleware seperti logger dan otelfiber diintegrasikan untuk logging dan tracing.
Route Definitions:

Mendefinisikan berbagai endpoint untuk API CRUD produk, seperti:
GET /product/:id: Mengambil data produk berdasarkan ID.
GET /product: Mengambil semua data produk.
POST /product: Menambahkan produk baru.
PUT /product/:id: Memperbarui data produk berdasarkan ID.
DELETE /product/:id: Menghapus produk berdasarkan ID.
Server Listening:

Aplikasi mendengarkan pada port 3000 untuk menerima request.
Komentar-komentar ini dapat membantu dalam memahami fungsi setiap bagian dari kode dan cara mereka berinteraksi satu sama lain untuk membentuk aplikasi CRUD API berbasis Hexagonal Architecture.
*/
