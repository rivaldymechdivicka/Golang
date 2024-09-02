package infrastructure

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
)

// AdapterMongo adalah struktur yang menyimpan konfigurasi dan client MongoDB
type AdapterMongo struct {
	MongoDSN string          // Data Source Name untuk koneksi MongoDB
	ctx      context.Context // Context untuk mengelola waktu hidup permintaan
	Client   *mongo.Client   // Client MongoDB yang digunakan untuk berinteraksi dengan database
	DB       string          // Nama database yang digunakan
}

// NewMongo adalah konstruktor untuk AdapterMongo. Ini membuat instance baru dari AdapterMongo dengan
// DSN dan nama database yang diberikan.
func NewMongo(ctx context.Context, mongoDSN string, db string) *AdapterMongo {
	return &AdapterMongo{ctx: ctx, MongoDSN: mongoDSN, DB: db}
}

// Connect menghubungkan ke MongoDB menggunakan opsi yang telah ditentukan. Ini juga mengaktifkan
// monitoring OpenTelemetry untuk MongoDB.
func (m *AdapterMongo) Connect() *AdapterMongo {
	// Membuat objek options.Client untuk konfigurasi client MongoDB
	clientOptions := options.Client()

	// Menambahkan monitor OpenTelemetry untuk pemantauan
	clientOptions.Monitor = otelmongo.NewMonitor()

	// Mengatur URI koneksi MongoDB dari DSN yang diberikan
	clientOptions.ApplyURI(m.MongoDSN)

	// Mencoba untuk terhubung ke MongoDB dengan menggunakan clientOptions yang telah disiapkan
	client, err := mongo.Connect(m.ctx, clientOptions)
	if err != nil {
		// Jika terjadi kesalahan saat menghubungkan, log kesalahan tersebut
		slog.ErrorContext(m.ctx, "error connect to mongo", slog.Any("err", err))
	}

	// Jika koneksi berhasil, log informasi bahwa MongoDB terhubung
	slog.InfoContext(m.ctx, "Mongodb connected.")

	// Mengembalikan instance AdapterMongo yang baru dengan client dan nama database yang diatur
	return &AdapterMongo{
		Client: client,
		DB:     m.DB,
	}
}

/*
Penjelasan:
Package infrastructure: Ini adalah package Go yang biasanya digunakan untuk mengatur infrastruktur atau layer penyimpanan data dari aplikasi.

Struct AdapterMongo: Struct ini mendefinisikan konfigurasi dan client MongoDB:

MongoDSN: Data Source Name untuk menghubungkan ke MongoDB.
ctx: Context yang digunakan untuk mengelola waktu hidup dan pembatalan operasi.
Client: Instance dari mongo.Client yang digunakan untuk berinteraksi dengan MongoDB.
DB: Nama database yang akan digunakan.
Function NewMongo: Konstruktor yang membuat instance baru dari AdapterMongo dengan parameter yang diberikan. Ini membantu dalam menginisialisasi objek dengan nilai awal yang sesuai.

Method Connect:

Membuat konfigurasi client MongoDB menggunakan options.Client().
Menambahkan otelmongo.NewMonitor() untuk memantau operasi MongoDB menggunakan OpenTelemetry.
Mengatur URI koneksi MongoDB dengan clientOptions.ApplyURI(m.MongoDSN).
Mencoba untuk menghubungkan ke MongoDB dan menangani error jika koneksi gagal.
Jika berhasil, menyimpan client yang terhubung dan mengembalikan instance AdapterMongo dengan client yang terhubung dan nama database yang diatur.
Catatan: Anda mungkin perlu menyesuaikan konfigurasi logging atau penggunaan context sesuai dengan kebutuhan aplikasi Anda.


*/
