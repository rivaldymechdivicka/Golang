package product

import "go.mongodb.org/mongo-driver/bson/primitive"

// Product merepresentasikan struktur data untuk entitas produk dalam sistem.
type Product struct {
	ID        primitive.ObjectID `json:"product_id,omitempty" bson:"_id,omitempty"` // ID unik produk yang dihasilkan oleh MongoDB
	Name      string             `json:"product_name" bson:"product_name"`          // Nama produk
	Stock     int64              `json:"stock" bson:"stock"`                        // Jumlah stok produk yang tersedia
	CreatedAt int64              `json:"created_at"`                                // Waktu (timestamp) saat produk dibuat
	UpdatedAt int64              `json:"updated_at"`                                // Waktu (timestamp) saat produk terakhir kali diperbarui
	DeletedAt int64              `json:"deleted_at"`                                // Waktu (timestamp) saat produk dihapus atau ditandai sebagai dihapus
}

// Filter digunakan untuk menentukan kriteria pencarian atau pemfilteran produk.
type Filter struct {
	Page      int    `json:"page"`      // Nomor halaman saat ini untuk pagination
	Limit     int    `json:"limit"`     // Jumlah maksimal item yang ditampilkan per halaman
	Latitude  string `json:"latitude"`  // Koordinat lintang untuk pencarian berbasis lokasi
	Longitude string `json:"longitude"` // Koordinat bujur untuk pencarian berbasis lokasi
	Keyword   string `json:"keyword"`   // Kata kunci untuk mencari produk berdasarkan nama atau atribut lainnya
}

/*

Penjelasan Fungsi Kode:
Deklarasi Package:

package product: Menyatakan bahwa file ini berada dalam package product. Package ini mengelompokkan kode yang berhubungan dengan produk agar lebih modular.
Import Statements:

import "go.mongodb.org/mongo-driver/bson/primitive": Mengimpor package yang digunakan untuk bekerja dengan tipe data ObjectID dari MongoDB, yang digunakan untuk mengelola ID produk.
Struct Product:

Struct Product digunakan untuk merepresentasikan data produk dalam sistem.
Field ID:
ID adalah ID unik produk yang dihasilkan oleh MongoDB dan disimpan dalam field _id.
Field Name:
Name menyimpan nama produk. Ketika data dikonversi menjadi JSON, field ini akan disebut product_name.
Field Stock:
Stock menyimpan jumlah stok yang tersedia untuk produk ini.
Field CreatedAt:
CreatedAt menyimpan waktu saat produk ini pertama kali dibuat dalam format UNIX timestamp.
Field UpdatedAt:
UpdatedAt menyimpan waktu saat produk ini terakhir kali diperbarui. Ini berguna untuk melacak perubahan yang dilakukan pada produk.
Field DeletedAt:
DeletedAt menyimpan waktu saat produk ini dihapus atau ditandai sebagai dihapus, juga dalam format UNIX timestamp.
Struct Filter:

Struct Filter digunakan untuk memfasilitasi pencarian atau pemfilteran produk berdasarkan kriteria tertentu.
Field Page:
Page digunakan untuk pagination, menunjukkan halaman mana yang sedang diakses.
Field Limit:
Limit menentukan jumlah maksimal produk yang akan dikembalikan per halaman.
Field Latitude:
Latitude digunakan untuk menyimpan koordinat lintang (latitude) untuk pencarian berbasis lokasi.
Field Longitude:
Longitude digunakan untuk menyimpan koordinat bujur (longitude) dalam pencarian berbasis lokasi.
Field Keyword:
Keyword digunakan untuk pencarian berdasarkan kata kunci, memungkinkan pengguna mencari produk berdasarkan nama atau atribut lain yang relevan.
Tujuan Komentar:
Komentar dalam bahasa Indonesia ini ditambahkan untuk menjelaskan tujuan dan fungsi dari setiap bagian kode. Komentar ini penting untuk memudahkan pemahaman kode, baik bagi Anda sendiri di masa depan atau bagi pengembang lain yang bekerja dengan kode ini.

*/
