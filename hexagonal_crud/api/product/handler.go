package product

import (
	"github.com/gofiber/fiber/v2"
)

// ProductHandler mendefinisikan sebuah interface dengan operasi-operasi
// yang harus diimplementasikan oleh handler yang spesifik,
// seperti handler HTTP, gRPC, atau protokol lainnya.
type ProductHandler interface {
	// Get mengambil satu entitas Product berdasarkan ID yang diterima dari konteks request.
	// Metode ini diharapkan mengarahkan permintaan ke service di domain untuk mengambil data produk.
	Get(ctx *fiber.Ctx)

	// Create membuat entitas Product baru dengan data yang diterima dari request.
	// Metode ini akan mengarahkan data ke service di domain untuk disimpan ke dalam database.
	Create(ctx *fiber.Ctx)

	// Update memperbarui entitas Product yang sudah ada dengan data yang diterima dari request.
	// Metode ini akan meneruskan data yang diperbarui ke service di domain untuk diupdate.
	Update(ctx *fiber.Ctx)

	// Delete menghapus entitas Product berdasarkan ID yang diterima dari konteks request.
	// Permintaan ini akan diteruskan ke service di domain yang bertanggung jawab untuk menghapus data produk.
	Delete(ctx *fiber.Ctx)

	// GetAll mengambil daftar semua entitas Product yang tersedia.
	// Service di domain akan mengembalikan daftar produk yang diminta.
	GetAll(ctx *fiber.Ctx)
}

/*
Penjelasan Tambahan:
Tujuan dari File handler.go:

File ini mendefinisikan interface ProductHandler yang berfungsi sebagai kontrak untuk berbagai operasi CRUD (Create, Read, Update, Delete) pada entitas Product. Interface ini memastikan bahwa setiap implementasi handler (misalnya untuk HTTP atau gRPC) akan menyediakan metode-metode dasar yang sama untuk mengelola entitas Product.
Manfaat dari Interface:

Penggunaan interface ini memberikan fleksibilitas dan modularitas pada arsitektur aplikasi, memungkinkan implementasi handler yang berbeda tanpa mengubah kode di bagian lain dari aplikasi, selama implementasi tersebut mematuhi kontrak yang ditetapkan oleh interface ini.
Konteks dalam Hexagonal Architecture:

Dalam konteks Hexagonal Architecture, file ini berada di lapisan "adapter", di mana ia menyediakan antarmuka antara dunia luar (misalnya HTTP request) dan logika bisnis inti yang berada di domain.
Komentar-komentar ini dirancang untuk memberikan pemahaman yang jelas mengenai peran setiap bagian kode di dalam file handler.go.
*/
