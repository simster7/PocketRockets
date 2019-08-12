package main

import (
	"flag"
	"fmt"
	"github.com/simster7/PocketRockets/go_backend/api/v1"
	"github.com/simster7/PocketRockets/go_backend/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 1234))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	v1.RegisterPokerServiceServer(grpcServer, &server.PokerServer{})
	reflection.Register(grpcServer)
	grpcServer.Serve(lis)
}
