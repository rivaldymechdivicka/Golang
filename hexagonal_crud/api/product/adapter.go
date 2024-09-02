package product

import (
	"CRUD_Hexagonal/domain/product"   // Mengimpor package domain untuk entitas Product
	_ "CRUD_Hexagonal/domain/product" // Mengimpor package domain (penggunaan implicit)
	"CRUD_Hexagonal/infrastructure"   // Mengimpor package infrastructure untuk tracing dan logging
	"CRUD_Hexagonal/utils"            // Mengimpor package utils untuk fungsi-fungsi utilitas seperti respon HTTP dan validasi
	"errors"                          // Mengimpor package errors untuk menangani error
	"net/http"                        // Mengimpor net/http untuk status code HTTP

	"github.com/gofiber/fiber/v2"                // Mengimpor Fiber untuk membuat handler HTTP
	"go.mongodb.org/mongo-driver/bson/primitive" // Mengimpor MongoDB primitive untuk menangani ObjectID
	"golang.org/x/exp/slog"                      // Mengimpor Slog untuk logging
)

// Struktur adapter yang menyimpan referensi ke service yang berhubungan dengan Product
type adapter struct {
	storeService product.ProductInterface
}

// NewStoreHandler adalah constructor untuk membuat instance dari adapter
// Fungsi ini menerima ProductInterface sebagai parameter, yang merupakan interface dari service di domain
func NewStoreHandler(storeService product.ProductInterface) *adapter {
	return &adapter{storeService: storeService}
}

// Fungsi Get adalah handler untuk endpoint GET /product/:id
// Fungsi ini mengambil satu product berdasarkan id yang diterima dari parameter URL
func (h *adapter) Get(ctx *fiber.Ctx) error {

	// Tracing dimulai untuk memantau eksekusi fungsi Get
	c, span := infrastructure.Tracer().Start(ctx.UserContext(), "api:store:Get")
	defer span.End()

	// Mengambil id dari URL parameter
	id := ctx.Params("id")
	// Memanggil service untuk mencari product berdasarkan id
	resp, err := h.storeService.Find(c, id)
	if err != nil {
		// Jika terjadi error, kembalikan response dengan status Bad Request
		utils.ResponseWithJSON(ctx, http.StatusBadRequest, *resp, err)
		return nil
	}
	// Jika berhasil, kembalikan response dengan status OK dan data product
	utils.ResponseWithJSON(ctx, http.StatusOK, resp, nil)
	return nil
}

// Fungsi GetAll adalah handler untuk endpoint GET /products
// Fungsi ini mengambil semua product berdasarkan filter yang diberikan melalui query parameters
func (h *adapter) GetAll(ctx *fiber.Ctx) error {
	// Tracing dimulai untuk memantau eksekusi fungsi GetAll
	c, span := infrastructure.Tracer().Start(ctx.UserContext(), "api:store:GetAll")
	defer span.End()

	// Membuat filter berdasarkan query parameters
	filter := product.Filter{
		Page:      ctx.QueryInt("page"),
		Limit:     ctx.QueryInt("limit"),
		Latitude:  ctx.Query("latitude"),
		Longitude: ctx.Query("longitude"),
		Keyword:   ctx.Query("keyword"),
	}
	// Memanggil service untuk mencari semua product berdasarkan filter
	p, pagination, err := h.storeService.FindAll(c, filter)
	if err != nil {
		// Jika terjadi error, kembalikan response dengan status Bad Request
		utils.ResponseWithJSON(ctx, http.StatusBadRequest, []*product.Product{}, err)
		return nil
	}
	// Jika berhasil, kembalikan response dengan status OK, data product, dan informasi pagination
	utils.ResponseWithJSON(ctx, http.StatusOK, p, nil, pagination)
	return nil
}

// Fungsi Create adalah handler untuk endpoint POST /product
// Fungsi ini membuat product baru berdasarkan data yang dikirim melalui request body
func (h *adapter) Create(ctx *fiber.Ctx) error {

	// Tracing dimulai untuk memantau eksekusi fungsi Create
	c, span := infrastructure.Tracer().Start(ctx.UserContext(), "api:store:Create")
	defer span.End()

	// Mem-parsing data dari request body ke dalam struct Product
	dataStore := &product.Product{}
	if err := ctx.BodyParser(&dataStore); err != nil {
		// Jika terjadi error saat parsing, kembalikan response dengan status Bad Request
		utils.ResponseWithJSON(ctx, http.StatusBadRequest, nil, err)
		return nil
	}

	// Melakukan validasi terhadap data product yang diterima
	errValidation := utils.Validate(dataStore)
	if errValidation != "" {
		// Jika validasi gagal, log error dan kembalikan response dengan status Unprocessable Entity
		slog.ErrorContext(c, "Failed to Validate api:product:Create", slog.Any("err ", errValidation))
		utils.ResponseWithJSON(ctx, http.StatusUnprocessableEntity, nil, errors.New(errValidation))
		return nil
	}

	// Memanggil service untuk menyimpan product baru
	resp, err := h.storeService.Store(c, dataStore)
	if err != nil {
		// Jika terjadi error saat penyimpanan, kembalikan response dengan status Bad Request
		utils.ResponseWithJSON(ctx, http.StatusBadRequest, nil, err)
		return nil
	}
	// Jika berhasil, kembalikan response dengan status OK dan data product yang disimpan
	utils.ResponseWithJSON(ctx, http.StatusOK, resp, nil)
	return nil

}

// Fungsi Update adalah handler untuk endpoint PUT /product/:id
// Fungsi ini memperbarui product yang ada berdasarkan id dan data baru yang dikirim melalui request body
func (h adapter) Update(ctx *fiber.Ctx) error {
	// Tracing dimulai untuk memantau eksekusi fungsi Update
	c, span := infrastructure.Tracer().Start(ctx.UserContext(), "api:store:Update")
	defer span.End()

	// Mengambil id dari URL parameter
	paramsID := ctx.Params("id")

	// Mem-parsing data dari request body ke dalam struct Product
	dataStore := &product.Product{}
	if err := ctx.BodyParser(&dataStore); err != nil {
		// Jika terjadi error saat parsing, kembalikan response dengan status Bad Request
		utils.ResponseWithJSON(ctx, http.StatusBadRequest, nil, err)
		return nil
	}

	// Melakukan validasi terhadap data product yang diterima
	errValidation := utils.Validate(dataStore)
	if errValidation != "" {
		// Jika validasi gagal, log error dan kembalikan response dengan status Unprocessable Entity
		slog.ErrorContext(c, "Failed to Validate api:product:Create", slog.Any("err ", errValidation))
		utils.ResponseWithJSON(ctx, http.StatusUnprocessableEntity, nil, errors.New(errValidation))
		return nil
	}

	// Memparsing id dari string ke ObjectID (untuk digunakan di MongoDB)
	id, errObjectID := primitive.ObjectIDFromHex(paramsID)
	if errObjectID != nil {
		// Jika terjadi error saat parsing ObjectID, log error dan kembalikan response dengan status Unprocessable Entity
		slog.ErrorContext(c, "Failed to Validate api:product:Create", slog.Any("err ", errObjectID))
		utils.ResponseWithJSON(ctx, http.StatusUnprocessableEntity, nil, errObjectID)
		return nil
	}
	// Mengatur id pada data product
	dataStore.ID = id

	// Memanggil service untuk memperbarui product
	err := h.storeService.Update(c, dataStore)
	if err != nil {
		// Jika terjadi error saat update, kembalikan response dengan status Bad Request
		utils.ResponseWithJSON(ctx, http.StatusBadRequest, nil, err)
		return nil
	}

	// Jika berhasil, kembalikan response dengan status OK dan data product yang diperbarui
	utils.ResponseWithJSON(ctx, http.StatusOK, dataStore, nil)
	return nil
}

// Fungsi Delete adalah handler untuk endpoint DELETE /product/:id
// Fungsi ini menghapus product berdasarkan id yang diterima dari parameter URL
func (h *adapter) Delete(ctx *fiber.Ctx) error {
	// Tracing dimulai untuk memantau eksekusi fungsi Delete
	c, span := infrastructure.Tracer().Start(ctx.UserContext(), "api:store:DeleteByID")
	defer span.End()

	// Mengambil id dari URL parameter
	id := ctx.Params("id")
	// Memanggil service untuk menghapus product berdasarkan id
	err := h.storeService.DeleteById(c, id)
	if err != nil {
		// Jika terjadi error saat penghapusan, kembalikan response dengan status Internal Server Error
		return ctx.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	// Jika berhasil, kembalikan response dengan status OK dan pesan sukses
	return ctx.Status(200).JSON(fiber.Map{"message": "Deleted successfully"})
}

/*
Berikut adalah penjelasan tambahan mengenai kode yang telah diberikan:

Struktur Kode dan Tujuan:
Paket yang Diimpor:

product: Ini adalah paket domain yang mendefinisikan entitas Product dan logika bisnis terkait. Paket ini berisi interface ProductInterface yang digunakan oleh lapisan adapter ini untuk berinteraksi dengan domain.
infrastructure: Paket ini mengandung infrastruktur aplikasi seperti tracing dan logging, yang digunakan untuk memantau dan mencatat aktivitas dalam aplikasi.
utils: Ini adalah paket yang berisi fungsi-fungsi utilitas yang digunakan untuk menangani respons HTTP, validasi, dan operasi umum lainnya.
errors: Digunakan untuk membuat dan menangani error dalam kode.
net/http: Paket ini menyediakan konstanta status kode HTTP yang digunakan dalam respons HTTP.
fiber: Framework HTTP yang digunakan untuk membangun API yang efisien dan ringan.
primitive: Bagian dari MongoDB driver yang digunakan untuk menangani ObjectID, tipe data unik yang digunakan sebagai kunci utama dalam MongoDB.
slog: Paket untuk logging, yang memungkinkan pencatatan aktivitas dalam aplikasi, terutama saat terjadi error atau eksekusi penting.
Struktur adapter:

Struktur adapter bertindak sebagai lapisan yang menghubungkan HTTP API dengan logika bisnis di domain. Ini menyimpan referensi ke service ProductInterface, yang merupakan interface domain untuk mengelola entitas Product.
Fungsi NewStoreHandler:

Fungsi ini adalah constructor yang membuat instance baru dari adapter dengan menerima service domain sebagai parameter. Fungsi ini menginisialisasi adapter dengan service domain yang akan digunakan untuk menangani operasi CRUD.
Fungsi CRUD:

Get: Mengambil satu entitas Product berdasarkan ID dari URL parameter.
GetAll: Mengambil semua entitas Product berdasarkan filter yang diterima dari query parameters.
Create: Membuat entitas Product baru berdasarkan data yang dikirim melalui request body.
Update: Memperbarui entitas Product yang ada berdasarkan ID dari URL parameter dan data baru yang dikirim melalui request body.
Delete: Menghapus entitas Product berdasarkan ID yang diterima dari URL parameter.
Tracing dan Logging:

Setiap fungsi menggunakan tracing yang dimulai dengan infrastructure.Tracer().Start() untuk memantau eksekusi fungsi tersebut. Tracing ini berguna untuk melacak alur eksekusi dalam aplikasi dan membantu dalam debugging.
Logging dilakukan dengan slog.ErrorContext() untuk mencatat error yang terjadi selama proses eksekusi, terutama saat validasi atau operasi lainnya gagal.
Validasi dan Error Handling:

Data yang diterima dari request body divalidasi menggunakan fungsi utils.Validate(). Jika validasi gagal, aplikasi akan mengembalikan status HTTP Unprocessable Entity (422).
Setiap error yang terjadi selama eksekusi fungsi (misalnya, parsing ID, validasi, atau operasi database) akan ditangani dengan mengembalikan respons yang sesuai, misalnya Bad Request (400) atau Internal Server Error (500).
Pentingnya Lapisan Adapter:
Lapisan adapter ini berfungsi sebagai jembatan antara dunia luar (seperti HTTP API) dan logika bisnis inti yang ada di domain. Ini memastikan bahwa segala interaksi dari klien (misalnya, browser, aplikasi mobile, atau layanan lain) diproses secara konsisten dan sesuai dengan aturan bisnis yang telah ditentukan. Dalam arsitektur heksagonal (Hexagonal Architecture), lapisan adapter adalah bagian penting yang memisahkan logika bisnis dari detail implementasi teknis seperti HTTP, sehingga memudahkan pemeliharaan, pengujian, dan pengembangan berkelanjutan.
*/
