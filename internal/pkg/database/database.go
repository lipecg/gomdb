package database

import (
	"context"
	"gomdb/internal/pkg/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dbConnString = "mongodb://gomdb-root:8lURb24nnHE8Kht3@10.0.0.126:27017/?retryWrites=true&w=majority"

func getMongoCollection(ms mongoStore) *mongo.Collection {
	return ms.Client.Database("gomdb").Collection("movies")
}

type mongoStore struct {
	Client *mongo.Client
}

func NewMongoStore() (domain.MovieDB, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dbConnString))
	if err != nil {
		return nil, err
	}
	return mongoStore{Client: client}, nil
}

func (ms mongoStore) Get(id int) (*domain.Movie, error) {
	col := getMongoCollection(ms)
	filter := bson.M{"id": id}
	movie := domain.Movie{}
	err := col.FindOne(context.TODO(), filter).Decode(&movie)
	if err != nil {
		return nil, err
	}
	return &movie, nil
}

// TODO: This won't work because of the query
func (ms mongoStore) List(query string) ([]*domain.Movie, error) {
	col := getMongoCollection(ms)
	filter := bson.M{"query": query}
	pets := []*domain.Movie{}
	cursor, err := col.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	if err := cursor.All(context.TODO(), &pets); err != nil {
		return nil, err
	}
	return pets, nil
}

func (ms mongoStore) Upsert(movie *domain.Movie) error {
	col := getMongoCollection(ms)

	filter := bson.M{"id": movie.ID}
	options := options.Replace().SetUpsert(true)

	result, err := col.ReplaceOne(context.TODO(), filter, movie, options)
	if result.UpsertedCount > 0 && result.UpsertedID != nil {
		movie.ObjectId = result.UpsertedID
	}

	return err
}
