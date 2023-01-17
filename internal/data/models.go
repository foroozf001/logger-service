// Package data provides an interface for logs
package data

import (
	"context"
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

// New is the function used to create an instance of the data package
func New(mongo *mongo.Client) Models {
	client = mongo

	return Models{
		LogItem: LogItem{},
	}
}

// Models describes all data models available to the application
type Models struct {
	LogItem LogItem
}

// LogItem describes a log object
type LogItem struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string    `bson:"name" json:"name"`
	Data      string    `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// Insert puts a document in a mongodb collection
func (l *LogItem) Insert(entry LogItem) error {
	databaseEnv, defined := os.LookupEnv("MONGO_INITDB_DATABASE")
	if !defined {
		return errors.New("mongo database error")
	}
	collectionEnv, defined := os.LookupEnv("MONGO_INITDB_COLLECTION")
	if !defined {
		return errors.New("mongo collection error")
	}
	ttlEnv, err := strconv.ParseInt(os.Getenv("MONGO_INITDB_TTL"), 10, 32)
	if err != nil {
		return errors.New("mongo ttl error")
	}

	collection := client.Database(databaseEnv).Collection(collectionEnv)

	// create index to expire documents
	index := mongo.IndexModel{
		Keys:    bson.D{{Key: "created_at", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(int32(ttlEnv) * 24 * 60 * 60), // convert days to seconds
	}

	_, err = collection.Indexes().CreateOne(context.Background(), index)
	if err != nil {
		return err
	}

	_, err = collection.InsertOne(context.TODO(), LogItem{
		Name:      entry.Name,
		Data:      entry.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return err
	}

	log.Println("logged event", entry.Name)

	return nil
}
