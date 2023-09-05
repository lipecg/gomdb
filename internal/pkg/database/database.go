package database

import (
	"context"
	"gomdb/internal/pkg/domain"
	"gomdb/internal/pkg/logging"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dbConnString = "mongodb://gomdb-root:8lURb24nnHE8Kht3@10.0.0.126:27017/?retryWrites=true&w=majority"

var client *mongo.Client
var db *mongo.Database

func init() {
	var err error
	client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(dbConnString))
	if err != nil {
		logging.Panic(err.Error())
	}
	db = client.Database("gomdb")
}

func getMongoCollection(collection string) *mongo.Collection {
	return db.Collection(collection)
}

func InsertOne(entity *domain.Entity, collection string) (*mongo.InsertOneResult, error) {

	col := getMongoCollection(collection)

	result, err := col.InsertOne(context.TODO(), entity)

	return result, err
}

func InsertMany(entities []interface{}, collection string) (*mongo.InsertManyResult, error) {

	col := getMongoCollection(collection)

	result, err := col.InsertMany(context.TODO(), entities)

	return result, err
}

func Upsert(entity *domain.Entity, collection string) (*mongo.UpdateResult, error) {

	entity.Updated = time.Now()

	col := getMongoCollection(collection)

	filter := bson.M{"id": entity.ID}
	options := options.Replace().SetUpsert(true)

	result, err := col.ReplaceOne(context.TODO(), filter, entity, options)

	return result, err
}

func Get(entity *domain.Entity, collection string) error {

	col := getMongoCollection(collection)

	filter := bson.M{"id": entity.ID}
	options := options.FindOne()
	options.Sort = bson.M{"updated": -1}

	err := col.FindOne(context.Background(), filter, options).Decode(entity)

	return err
}

type ListOptions struct {
	Page        int    `default:"1"`
	PageSize    int    `default:"10"`
	SearchText  string `default:""`
	SearchField string `default:""`
}

type ListResult struct {
	TotalRecords int             `json:"total_results"`
	TotalPages   int             `json:"total_pages"`
	CurrentPage  int             `json:"current_page"`
	CountResults int             `json:"count_results"`
	PageSize     int             `json:"page_size"`
	List         []domain.Entity `json:"list"`
}

func List(collection string, listResult *ListResult, listOptions ListOptions) error {

	col := getMongoCollection(collection)

	filter := bson.M{}
	options := options.Find()
	options.Sort = bson.M{"updated": -1}
	options.SetLimit(int64(listOptions.PageSize))
	options.SetSkip(int64(listOptions.PageSize * (listOptions.Page - 1)))

	if listOptions.SearchText != "" && listOptions.SearchField != "" {
		filter[listOptions.SearchField] = bson.M{"$regex": listOptions.SearchText, "$options": "i"}
	}

	// get the total number of documents
	total, err := col.CountDocuments(context.Background(), filter)
	if err != nil {
		return err
	}

	listResult.TotalRecords = int(total)
	listResult.TotalPages = int(math.Ceil(float64(total) / float64(listOptions.PageSize)))
	listResult.CurrentPage = listOptions.Page
	listResult.PageSize = listOptions.PageSize

	//cursor, err := col.Find(context.Background(), filter, options)

	pipeline := []bson.M{
		{"$match": filter},
		{"$sort": bson.M{"id": 1, "updated": -1}},
		{"$group": bson.M{"_id": "$id", "doc": bson.M{"$first": "$$ROOT"}}},
		{"$sort": bson.M{"_id": 1}},
		{"$skip": int64(listOptions.PageSize * (listOptions.Page - 1))},
		{"$limit": int64(listOptions.PageSize)},
	}

	cursor, err := col.Aggregate(context.Background(), pipeline)
	if err != nil {
		return err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var result struct {
			ID  int           `bson:"id"`
			Doc domain.Entity `bson:"doc"`
		}
		err := cursor.Decode(&result)
		if err != nil {
			return err
		}

		listResult.List = append(listResult.List, result.Doc)
		listResult.CountResults++
	}

	return nil
}

// db.movies.aggregate(
//    [
// 	{ $match: { id: 346698 } },
//      { $sort: { id: 1, updated: -1 } },
//      {
//        $group:
//          {
//            _id: "$id",
//            first: { $first: "$$ROOT" }
//          }
//      }
//    ]
// )
