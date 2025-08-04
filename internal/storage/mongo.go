package storage

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStore struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoStore(uri, dbName, collectionName string) (*MongoStore, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}

	coll := client.Database(dbName).Collection(collectionName)

	return &MongoStore{
		client:     client,
		collection: coll,
	}, nil
}

func (m *MongoStore) Save(code, url string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	doc := map[string]interface{}{
		"code": code,
		"url":  url,
	}

	_, err := m.collection.InsertOne(ctx, doc)
	return err
}

func (m *MongoStore) Get(code string) (string, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result struct {
		URL string `bson:"url"`
	}

	filter := map[string]interface{}{
		"code": code,
	}

	err := m.collection.FindOne(ctx, filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}

	return result.URL, true, nil
}
