package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/foroozf001/logger-service/internal/api"
	"github.com/foroozf001/logger-service/internal/data"
)

const (
	httpPort = "8080"
	grpcPort = "50051"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

var customCounter = promauto.NewCounter(prometheus.CounterOpts{
	Name: "custom_counter_total",
	Help: "The total number custom events",
})

func main() {
	mongoClient, err := api.ConnectMongo()
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

	// app := Config{
	// 	Models: data.New(client),
	// }

	log.Printf("listening on http port %s\n", httpPort)

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)
}
