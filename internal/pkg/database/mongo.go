package mongo

import (
	"context"
	"gomdb/cli/internal/pkg/models"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dbConnString = "mongodb://gomdb-root:8lURb24nnHE8Kht3@10.0.0.126:27017/?retryWrites=true&w=majority"

func getMovieFromDB(id int) models.Movie {

	moviesCollection := CNX.Database("gomdb").Collection("movies")

	movie := models.Movie{}

	filter := bson.D{primitive.E{Key: "id", Value: id}}

	err := moviesCollection.FindOne(context.TODO(), filter).Decode(&movie)

	if err != nil {
		if err != mongo.ErrNoDocuments {
			log.Panic(err)
		}
	}

	return movie
}

func UpdateMovieDB(movie *models.Movie) *mongo.UpdateResult {
	moviesCollection := CNX.Database("gomdb").Collection("movies")

	movie.Updated = time.Now()

	filter := bson.D{bson.E{Key: "id", Value: movie.ID}}
	doc, err := toDoc(movie)
	options := options.Replace().SetUpsert(true)

	if err != nil {
		log.Panic(err)
	}

	result, err := moviesCollection.ReplaceOne(context.TODO(), filter, doc, options)

	if err != nil {
		if err != mongo.ErrNoDocuments {
			log.Panic(err)
		}
	}

	return result
}

func InsertManyMovies(movies []interface{}) *mongo.InsertManyResult {

	moviesCollection := CNX.Database("gomdb").Collection("movies")

	result, err := moviesCollection.InsertMany(context.TODO(), movies)

	if err != nil {
		if err != mongo.ErrNoDocuments {
			log.Panic(err)
		}
	}

	return result
}

func toDoc(v interface{}) (doc *bson.D, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}

	err = bson.Unmarshal(data, &doc)
	return
}

func Connection() *mongo.Client {

	client, err := mongo.NewClient(options.Client().ApplyURI(dbConnString))

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)

	if err != nil {
		log.Fatal(err)
	}

	return client
}

var CNX = Connection()
