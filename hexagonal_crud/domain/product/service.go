package product

import (
	"CRUD_Hexagonal/utils"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ProductInterface mendefinisikan kontrak (interface) untuk layanan (service) produk.
type ProductInterface interface {
	// Find mencari produk berdasarkan ID dan mengembalikan produk jika ditemukan.
	Find(ctx context.Context, id string) (*Product, error)

	// Store menyimpan produk baru ke dalam database dan mengembalikan produk yang disimpan.
	Store(ctx context.Context, product *Product) (*Product, error)

	// Update memperbarui data produk yang ada di dalam database.
	Update(ctx context.Context, dataStore *Product) error

	// FindAll mencari semua produk berdasarkan filter yang diberikan dan mengembalikan daftar produk serta pagination.
	FindAll(ctx context.Context, filter Filter) ([]*Product, *utils.Pagination, error)

	// Delete menghapus produk dari database berdasarkan kode produk.
	Delete(ctx context.Context, code string) error

	// DeleteById menghapus produk dari database berdasarkan ID produk.
	DeleteById(ctx context.Context, id string) error
}

// Repository mendefinisikan kontrak (interface) untuk repository produk yang akan berinteraksi langsung dengan database.
type Repository interface {
	// Find mencari produk berdasarkan ID dan mengembalikan produk jika ditemukan.
	Find(ctx context.Context, id string) (*Product, error)

	// Store menyimpan produk baru ke dalam database dan mengembalikan ObjectID dari produk yang disimpan.
	Store(ctx context.Context, dataStore *Product) (primitive.ObjectID, error)

	// Update memperbarui data produk yang ada di dalam database.
	Update(ctx context.Context, dataStore *Product) error

	// FindAll mencari semua produk berdasarkan filter yang diberikan dan mengembalikan daftar produk serta pagination.
	FindAll(ctx context.Context, filter Filter) ([]*Product, *utils.Pagination, error)

	// Delete menghapus produk dari database berdasarkan kode produk.
	Delete(ctx context.Context, code string) error

	// DeleteById menghapus produk dari database berdasarkan ID produk.
	DeleteById(ctx context.Context, id string) error
}

/*
Penjelasan Fungsi Kode:
Deklarasi Package:

package product: Menyatakan bahwa file ini adalah bagian dari package product. Package ini mengelompokkan kode yang berhubungan dengan manajemen produk.
Import Statements:

import ("CRUD_Hexagonal/utils" "context" "go.mongodb.org/mongo-driver/bson/primitive"): Mengimpor package yang dibutuhkan seperti utils untuk pagination, context untuk manajemen konteks, dan primitive untuk menangani tipe data ObjectID dari MongoDB.
Interface ProductInterface:

ProductInterface adalah kontrak untuk layanan produk yang mendefinisikan fungsi-fungsi yang harus diimplementasikan oleh service yang mengelola produk.

Find:

Find(ctx context.Context, id string) (*Product, error): Fungsi ini bertanggung jawab untuk mencari produk berdasarkan id dan mengembalikan produk tersebut jika ditemukan.
Store:

Store(ctx context.Context, product *Product) (*Product, error): Fungsi ini digunakan untuk menyimpan produk baru ke dalam database dan mengembalikan produk yang baru disimpan.
Update:

Update(ctx context.Context, dataStore *Product) error: Fungsi ini memperbarui informasi produk yang ada di dalam database.
FindAll:

FindAll(ctx context.Context, filter Filter) ([]*Product, *utils.Pagination, error): Fungsi ini mencari semua produk berdasarkan filter yang diberikan dan mengembalikan daftar produk serta informasi pagination.
Delete:

Delete(ctx context.Context, code string) error: Fungsi ini menghapus produk dari database berdasarkan kode produk.
DeleteById:

DeleteById(ctx context.Context, id string) error: Fungsi ini menghapus produk dari database berdasarkan ID produk.
Interface Repository:

Repository adalah kontrak untuk repository produk yang berinteraksi langsung dengan database. Interface ini mendefinisikan fungsi-fungsi yang harus diimplementasikan oleh repository.

Find:

Find(ctx context.Context, id string) (*Product, error): Sama seperti di ProductInterface, fungsi ini mencari produk berdasarkan ID dan mengembalikan produk jika ditemukan.
Store:

Store(ctx context.Context, dataStore *Product) (primitive.ObjectID, error): Fungsi ini menyimpan produk baru ke dalam database dan mengembalikan ObjectID dari produk yang baru disimpan.
Update:

Update(ctx context.Context, dataStore *Product) error: Fungsi ini memperbarui data produk yang ada di dalam database.
FindAll:

FindAll(ctx context.Context, filter Filter) ([]*Product, *utils.Pagination, error): Fungsi ini mencari semua produk berdasarkan filter yang diberikan dan mengembalikan daftar produk serta pagination.
Delete:

Delete(ctx context.Context, code string) error: Fungsi ini menghapus produk dari database berdasarkan kode produk.
DeleteById:

DeleteById(ctx context.Context, id string) error: Fungsi ini menghapus produk dari database berdasarkan ID produk.
Tujuan Komentar:
Komentar dalam kode ini bertujuan untuk memberikan penjelasan tentang fungsi-fungsi dan interface yang didefinisikan dalam kode, serta bagaimana fungsi tersebut digunakan dalam konteks pengelolaan produk dalam aplikasi. Komentar ini sangat membantu untuk pemahaman dan pemeliharaan kode di masa depan.


*/
