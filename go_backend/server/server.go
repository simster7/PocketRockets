package server

import (
	"context"
	"flag"
	"fmt"
	"github.com/simster7/PocketRockets/go_backend/api/v1"
	"github.com/simster7/PocketRockets/go_backend/engine"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type PokerServer struct {
	Game engine.Game
	Players map[int32]engine.Player
}

func (s *PokerServer) GetPlayerState(context context.Context, request *v1.PlayerStateRequest) (*v1.PlayerState, error) {
	player := s.Players[request.PlayerId]
	return s.Game.GetPlayerState(&player), nil
}

func StartDevServer() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 1234))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	v1.RegisterPokerServiceServer(grpcServer, &PokerServer{})
	reflection.Register(grpcServer)
	grpcServer.Serve(lis)
}