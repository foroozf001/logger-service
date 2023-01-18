// Package api provides the web interface
package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/foroozf001/logger-service/internal/data"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	httpPort = "8080"
	grpcPort = "50051"
)

var client *mongo.Client

var totalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of HTTP requests processed, based on path",
	},
	[]string{"path"},
)

var responseStatus = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "response_status",
		Help: "Status of HTTP responses",
	},
	[]string{"status"},
)

var httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name: "http_response_time_seconds",
	Help: "Duration of HTTP requests",
}, []string{"path"})

type Config struct {
	Models data.Models
}

func init() {
	prometheus.Register(totalRequests)
	prometheus.Register(responseStatus)
	prometheus.Register(httpDuration)
}

// Run starts to listen on web and grpc ports
func Run() {
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

	app := Config{
		Models: data.New(client),
	}

	go app.GrpcListen(grpcPort)

	log.Printf("listening on http port %s\n", httpPort)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%s", httpPort),
		Handler:           app.Routes(),
		ReadHeaderTimeout: 2 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
