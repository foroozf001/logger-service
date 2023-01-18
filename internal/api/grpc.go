// Package api provides the web interface
package api

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/foroozf001/logger-service/internal/data"
	pb "github.com/foroozf001/logger-service/internal/proto/v1"
	"google.golang.org/grpc"
)

// LogServer is type used for writing to the log via gRPC
type LogServer struct {
	pb.UnimplementedLogServiceServer
	Models data.Models
}

// WriteLog writes the log after receiving a call from a gRPC client
func (l *LogServer) WriteLog(ctx context.Context, req *pb.LogRequest) (*pb.LogResponse, error) {
	input := req.GetLogEntry()

	LogItem := data.LogItem{
		Name: input.Name,
		Data: input.Data,
	}

	err := l.Models.LogItem.Insert(LogItem)
	if err != nil {
		res := &pb.LogResponse{Result: "failed"}
		return res, err
	}

	res := &pb.LogResponse{Result: "logged"}
	return res, nil
}

func (app *Config) GrpcListen(grpcPort string) {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen for tcp: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterLogServiceServer(s, &LogServer{Models: app.Models})

	log.Printf("listening on grpc port %s\n", grpcPort)
	if err := s.Serve(listen); err != nil {
		log.Fatalf("failed to listen for grpc: %v", err)
	}
}
