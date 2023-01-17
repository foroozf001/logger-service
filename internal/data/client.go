// Package data provides an interface for logs
package data

import (
	"context"
	"errors"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongo() (*mongo.Client, error) {
	usernameEnv, defined := os.LookupEnv("MONGO_INITDB_ROOT_USERNAME")
	if !defined {
		return nil, errors.New("mongo username error")
	}
	passwordEnv, defined := os.LookupEnv("MONGO_INITDB_ROOT_PASSWORD")
	if !defined {
		return nil, errors.New("mongo password error")
	}
	uriEnv, defined := os.LookupEnv("MONGO_SERVICE_URI")
	if !defined {
		return nil, errors.New("mongo uri error")
	}

	clientOptions := options.Client().ApplyURI(uriEnv)
	clientOptions.SetAuth(options.Credential{
		Username: usernameEnv,
		Password: passwordEnv,
	})

	connection, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	log.Printf("connected %s\n", uriEnv)

	return connection, nil
}
