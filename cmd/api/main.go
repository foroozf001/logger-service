package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/foroozf001/logger-service/internal/data"
	"github.com/prometheus/client_golang/prometheus"
)

type Config struct {
	Models data.Models
}

type Config2 struct {
	Models data.Models
}

var httpReqs = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "How many HTTP requests processed, partitioned by status code and http method.",
	},
	[]string{"code", "method"},
)

const (
	httpPort = "8080"
	grpcPort = "50051"
)

var client *mongo.Client

func main() {
	// prometheus.MustRegister(httpReqs)

	// mongoClient, err := ConnectMongo()
	// if err != nil {
	// 	log.Panic(err)
	// }

	// client = mongoClient

	// ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	// defer cancel()

	// defer func() {
	// 	if err = client.Disconnect(ctx); err != nil {
	// 		panic(err)
	// 	}
	// }()

	app := Config2{
		Models: data.New(client),
	}

	log.Printf("listening on http port %s\n", httpPort)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%s", httpPort),
		Handler:           app.routes(),
		ReadHeaderTimeout: 2 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

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
