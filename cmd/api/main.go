package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/foroozf001/logger-service/internal/api"
	"github.com/foroozf001/logger-service/internal/data"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	httpPort = "8080"
	grpcPort = "50051"
)

var client *mongo.Client

var httpReqs = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "How many HTTP requests processed, partitioned by status code and http method.",
	},
	[]string{"code", "method"},
)

func main() {
	prometheus.MustRegister(httpReqs)

	mongoClient, err := data.ConnectMongo()
	if err != nil {
		log.Panic(err)
	}

	client = mongoClient

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := api.Config{
		Models: data.New(client),
		HttpReqs: httpReqs,
	}

	log.Printf("listening on http port %s\n", httpPort)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%s", httpPort),
		Handler:           app.Routes(),
		ReadHeaderTimeout: 2 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
