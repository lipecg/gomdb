package database

import (
	"context"
	"errors"
	"fmt"
	"gomdb/internal/pkg/domain"
	"gomdb/internal/pkg/logging"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getMongoCollection(ms mongoStore, collection ...string) (*mongo.Collection, error) {
	if collection != nil {
		if collection[0] == "*domain.Movie" {
			return ms.Client.Database("gomdb").Collection("movies"), nil
		} else if collection[0] == "*domain.TVSeries" {
			return ms.Client.Database("gomdb").Collection("tvseries"), nil
		} else if collection[0] == "*domain.Person" {
			return ms.Client.Database("gomdb").Collection("people"), nil
		}
	}
	return nil, errors.New("Cannot get collection")
}

type mongoStore struct {
	Client *mongo.Client
}

func NewMongoStore(dbConnString string) (domain.EntityDB, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dbConnString))
	if err != nil {
		logging.Panic(err.Error())
	}
	return mongoStore{Client: client}, nil
}

func (ms mongoStore) Get(id int) (*interface{}, error) {
	//TODO: handle error properly
	//TODO: get proper type
	col, _ := getMongoCollection(ms, "*domain.Movie")
	filter := bson.M{"id": id}
	var entity interface{}
	err := col.FindOne(context.TODO(), filter).Decode(&entity)
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// TODO: This won't work because of the query
// TODO: Need to get collection name from .... Query?
func (ms mongoStore) List(query string) ([]*interface{}, error) {
	//TODO: handle error properly
	col, _ := getMongoCollection(ms, "*domain.Movie")
	filter := bson.M{"query": query}
	var entities []*interface{}
	cursor, err := col.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	if err := cursor.All(context.TODO(), &entities); err != nil {
		return nil, err
	}
	return entities, nil
}

func (ms mongoStore) Upsert(entity *interface{}) error {
	entityType := reflect.TypeOf(*entity).String()
	col, err := getMongoCollection(ms, entityType)
	if err != nil {
		return err
	}

	filter := bson.M{"id": reflect.ValueOf(*entity).Elem().FieldByName("ID").Int()}
	options := options.Replace().SetUpsert(true)

	result, err := col.ReplaceOne(context.TODO(), filter, entity, options)

	fmt.Print(result)

	return err
}
