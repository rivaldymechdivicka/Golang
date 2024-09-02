package product

import (
	"CRUD_Hexagonal/domain/product"
	"CRUD_Hexagonal/infrastructure"
	"CRUD_Hexagonal/utils"
	"context"
	"time"
)

// adapter adalah struct yang mengimplementasikan interface ProductInterface
// untuk berinteraksi dengan repository produk (store).
type adapter struct {
	storeRepo product.Repository
}

// NewStoreService adalah constructor yang digunakan untuk membuat instance baru dari adapter
// dan mengembalikannya sebagai implementasi ProductInterface.
func NewStoreService(storeRepo product.Repository) product.ProductInterface {
	return &adapter{storeRepo: storeRepo}
}

// Find mencari produk (store) berdasarkan ID yang diberikan.
func (a adapter) Find(ctx context.Context, id string) (*product.Product, error) {
	// Memulai tracing untuk fungsi Find
	ctx, span := infrastructure.Tracer().Start(ctx, "service:store:Find")
	defer span.End()

	return a.storeRepo.Find(ctx, id)
}

// Store menyimpan produk (store) baru ke dalam repository dan mengatur waktu pembuatan.
func (a adapter) Store(ctx context.Context, product *product.Product) (*product.Product, error) {
	// Memulai tracing untuk fungsi Store
	ctx, span := infrastructure.Tracer().Start(ctx, "service:store:Store")
	defer span.End()

	// Mengatur waktu pembuatan produk
	product.CreatedAt = time.Now().UTC().Unix()

	// Menyimpan produk ke dalam repository
	insertID, err := a.storeRepo.Store(ctx, product)

	// Mengatur ID produk dengan ID yang baru disisipkan
	product.ID = insertID

	return product, err
}

// Update memperbarui data produk (store) yang ada di dalam repository
// dan mengatur waktu pembaruan.
func (a adapter) Update(ctx context.Context, store *product.Product) error {
	// Memulai tracing untuk fungsi Update
	ctx, span := infrastructure.Tracer().Start(ctx, "service:store:Update")
	defer span.End()

	// Mengatur waktu pembaruan produk
	store.UpdatedAt = time.Now().UTC().Unix()

	return a.storeRepo.Update(ctx, store)
}

// FindAll mencari semua produk (store) dengan filter tertentu dan mendukung pagination.
func (a adapter) FindAll(ctx context.Context, filter product.Filter) ([]*product.Product, *utils.Pagination, error) {
	// Memulai tracing untuk fungsi FindAll
	ctx, span := infrastructure.Tracer().Start(ctx, "service:store:FindAll")
	defer span.End()

	// Mencari semua produk yang sesuai dengan filter dan mengembalikan hasil serta pagination
	res, pagination, err := a.storeRepo.FindAll(ctx, filter)

	return res, pagination, err
}

// Delete adalah fungsi yang belum diimplementasikan untuk menghapus produk (store)
// berdasarkan kode tertentu.
func (a adapter) Delete(ctx context.Context, code string) error {
	//TODO implement me
	panic("implement me")
}

// DeleteById menghapus produk (store) dari repository berdasarkan ID yang diberikan.
func (a adapter) DeleteById(ctx context.Context, id string) error {
	// Memulai tracing untuk fungsi DeleteById
	ctx, span := infrastructure.Tracer().Start(ctx, "service:store:DeleteByID")
	defer span.End()

	err := a.storeRepo.DeleteById(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

/*
Penjelasan Fungsi dan Komentar dalam Kode:
Package product:

Package ini digunakan untuk mengimplementasikan layanan yang berinteraksi dengan repository produk (store).
Struct adapter:

adapter adalah implementasi dari interface ProductInterface. Struct ini menggunakan storeRepo, yang merupakan instance dari product.Repository, untuk berinteraksi dengan repository produk.
Fungsi NewStoreService:

Fungsi ini adalah constructor untuk membuat instance baru dari adapter dan mengembalikannya sebagai implementasi dari ProductInterface. Ini memungkinkan layanan produk untuk digunakan di seluruh aplikasi.
Fungsi Find:

Fungsi ini mencari produk berdasarkan ID dan mengembalikannya. Fungsi ini juga memulai tracing untuk memantau kinerja dan masalah.
Fungsi Store:

Fungsi ini menyimpan produk baru ke dalam repository dan mengatur waktu pembuatan produk. Setelah produk disimpan, fungsi ini mengatur ID produk dengan ID yang baru disisipkan.
Fungsi Update:

Fungsi ini memperbarui produk yang ada di dalam repository dan mengatur waktu pembaruan produk. Ini memanfaatkan tracing untuk memantau proses.
Fungsi FindAll:

Fungsi ini mencari semua produk dengan filter tertentu dan mendukung pagination. Ini mengembalikan hasil pencarian serta informasi pagination.
Fungsi Delete:

Fungsi ini adalah placeholder yang belum diimplementasikan. Fungsi ini dimaksudkan untuk menghapus produk berdasarkan kode tertentu.
Fungsi DeleteById:

Fungsi ini menghapus produk dari repository berdasarkan ID yang diberikan. Fungsi ini juga memulai tracing untuk memantau kinerja dan masalah.
Dengan penjelasan dan komentar ini, diharapkan kode lebih mudah dipahami dan dimengerti fungsinya dalam konteks aplikasi.
*/
