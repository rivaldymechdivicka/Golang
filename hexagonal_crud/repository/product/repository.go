package product

import (
	"CRUD_Hexagonal/domain/product"
	"CRUD_Hexagonal/infrastructure"
	"CRUD_Hexagonal/utils"
	"context"
	"errors"
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/exp/slog"
)

// storeRepository adalah struct yang mengimplementasikan interface Repository
// untuk berinteraksi dengan koleksi MongoDB yang menyimpan data produk (store).
type storeRepository struct {
	client     *mongo.Client
	db         string
	collection string
}

// NewstoreRepository adalah constructor yang digunakan untuk membuat instance baru dari storeRepository.
func NewstoreRepository(client *mongo.Client, db string, collection string) product.Repository {
	return &storeRepository{
		client:     client,
		db:         db,
		collection: collection,
	}
}

// Find berfungsi untuk mencari satu produk (store) berdasarkan ID yang diberikan.
func (r *storeRepository) Find(ctx context.Context, id string) (*product.Product, error) {
	// Mulai tracing untuk fungsi Find
	ctx, span := infrastructure.Tracer().Start(ctx, "repository:store:Find")
	defer span.End()

	var storeData product.Product

	// Mengambil koleksi yang dituju dalam database
	collection := r.client.Database(r.db).Collection(r.collection)

	// Mengonversi ID dari string ke ObjectID MongoDB
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		slog.ErrorContext(ctx, "Collection Error", slog.Any("err ", err))
		return nil, err
	}

	// Membuat filter untuk pencarian berdasarkan ID
	filter := bson.D{{"_id", objectId}}
	err = collection.FindOne(ctx, filter).Decode(&storeData)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("error Finding a store")
		}
		slog.ErrorContext(ctx, "Collection Error", slog.Any("err ", err))
	}
	return &storeData, nil
}

// FindAll berfungsi untuk mencari semua produk (store) dengan filter tertentu dan mendukung pagination.
func (r *storeRepository) FindAll(ctx context.Context, filter product.Filter) ([]*product.Product, *utils.Pagination, error) {
	// Mulai tracing untuk fungsi FindAll
	ctx, span := infrastructure.Tracer().Start(ctx, "repository:store:FindAll")
	defer span.End()

	collection := r.client.Database(r.db).Collection(r.collection)

	// Pengaturan pagination
	var currentPage, limit int
	if filter.Limit <= 0 || filter.Page <= 0 {
		currentPage, limit = 1, 10
	} else {
		currentPage, limit = filter.Page, filter.Limit
	}
	skip := (currentPage - 1) * limit

	// Menghitung total dokumen untuk pagination
	totalDocuments, err := collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return nil, nil, err
	}

	// Membuat struktur pagination
	pagination := utils.Pagination{
		Total:       int(totalDocuments),
		Limit:       limit,
		CurrentPage: currentPage,
	}

	// Menerapkan pagination pada query
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))

	// Membuat filter untuk latitude dan longitude
	bsonFilter := bson.D{}

	if filter.Latitude != "" && filter.Longitude != "" {
		bsonFilter = append(bsonFilter, bson.E{
			Key:   "address.geo.latitude",
			Value: bson.D{{Key: "$eq", Value: filter.Latitude}},
		}, bson.E{
			Key:   "address.geo.longitude",
			Value: bson.D{{Key: "$eq", Value: filter.Longitude}},
		})
	}

	// Membuat filter untuk keyword pencarian
	if filter.Keyword != "" {
		bsonFilter = append(bsonFilter, bson.E{
			Key:   "name",
			Value: bson.D{{Key: "$regex", Value: primitive.Regex{Pattern: filter.Keyword, Options: "i"}}},
		})
	}

	var stores []*product.Product

	// Menjalankan query untuk menemukan data yang sesuai dengan filter dan opsi yang diterapkan
	cur, err := collection.Find(ctx, bsonFilter, findOptions)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil, errors.New("error finding store")
		}
		slog.ErrorContext(ctx, "Collection Error", slog.Any("err ", err))
		return nil, nil, err
	}
	defer cur.Close(ctx)

	// Memasukkan hasil query ke dalam slice stores
	for cur.Next(ctx) {
		var elem product.Product
		if err := cur.Decode(&elem); err != nil {
			slog.ErrorContext(ctx, "Collection Error", slog.Any("err ", err))
			continue
		}
		stores = append(stores, &elem)
	}
	return stores, &pagination, nil
}

// Store berfungsi untuk menyimpan data produk (store) baru ke dalam database.
func (r *storeRepository) Store(ctx context.Context, dataStore *product.Product) (primitive.ObjectID, error) {
	// Mulai tracing untuk fungsi Store
	ctx, span := infrastructure.Tracer().Start(ctx, "repository:store:Store")
	defer span.End()

	collection := r.client.Database(r.db).Collection(r.collection)

	// Menyisipkan data baru ke dalam koleksi
	doInsert, err := collection.InsertOne(ctx, dataStore)
	if err != nil {
		slog.ErrorContext(ctx, "Error writing to repository", slog.Any("err ", err))
		return primitive.ObjectID{}, errors.New("error writing to repository")
	}

	// Mengembalikan ID dari dokumen yang baru saja disimpan
	return doInsert.InsertedID.(primitive.ObjectID), nil
}

// Update berfungsi untuk memperbarui data produk (store) yang sudah ada di dalam database.
func (r *storeRepository) Update(ctx context.Context, dataStore *product.Product) error {
	// Mulai tracing untuk fungsi Update
	ctx, span := infrastructure.Tracer().Start(ctx, "repository:store:Update")
	defer span.End()

	updatedStore := bson.D{}

	// Menggunakan refleksi untuk memeriksa setiap field dalam struct dataStore
	values := reflect.ValueOf(*dataStore)
	types := values.Type()
	for i := 0; i < values.NumField(); i++ {
		// Jika field bukan ID dan bukan field kosong, maka tambahkan ke updatedStore
		if types.Field(i).Name != "ID" && !utils.IsEmptyStruct(values.Field(i)) {
			updatedStore = append(updatedStore, primitive.E{Key: types.Field(i).Tag.Get("json"), Value: values.Field(i).Interface()})
		}
	}

	collection := r.client.Database(r.db).Collection(r.collection)

	// Melakukan update pada dokumen berdasarkan ID
	_, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": dataStore.ID},
		bson.D{
			{Key: "$set", Value: updatedStore},
		},
	)
	if err != nil {
		fmt.Println(err, "err")
		return err
	}

	return nil
}

// DeleteById berfungsi untuk menghapus produk (store) dari database berdasarkan ID yang diberikan.
func (r *storeRepository) DeleteById(ctx context.Context, id string) error {
	// Mulai tracing untuk fungsi DeleteById
	ctx, span := infrastructure.Tracer().Start(ctx, "repository:store:DeleteById")
	defer span.End()

	// Periksa apakah ID kosong
	if id == "" {
		return errors.New("ID is empty")
	}

	collection := r.client.Database(r.db).Collection(r.collection)

	// Mengonversi ID dari string ke ObjectID MongoDB
	objectID, err2 := primitive.ObjectIDFromHex(id)
	if err2 != nil {
		return err2
	}

	// Memeriksa apakah dokumen dengan ID tersebut ada
	result := collection.FindOne(ctx, bson.M{"_id": objectID})

	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return errors.New("store not found")
		}
		return result.Err()
	}

	// Menghapus dokumen dari koleksi
	_, err := collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	return nil
}

// Delete berfungsi untuk menghapus produk (store) dari database berdasarkan kode tertentu.
// Fungsi ini belum diimplementasikan.
func (r *storeRepository) Delete(ctx context.Context, code string) error {
	//TODO implement me
	panic("implement me")
}

/*
Penjelasan Fungsi dan Komentar dalam Kode:
Package product:

Package ini digunakan untuk mengimplementasikan repository yang berinteraksi dengan database MongoDB, khususnya untuk menyimpan, mengambil, memperbarui, dan menghapus data produk (store).
Struct storeRepository:

storeRepository adalah implementasi dari interface product.Repository yang bertugas melakukan operasi CRUD pada koleksi MongoDB.
Fungsi NewstoreRepository:

Fungsi ini adalah constructor untuk membuat instance baru dari storeRepository. Parameter yang diterima adalah client (

*/
