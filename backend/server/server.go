package server

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/simster7/PocketRockets/backend/api/v1"
	"github.com/simster7/PocketRockets/backend/engine"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type PokerServer struct {
	Games   map[int32]*engine.Game
	Players map[int32]*engine.Player
}

func NewPokerServer() PokerServer {
	return PokerServer{
		Games:   make(map[int32]*engine.Game),
		Players: make(map[int32]*engine.Player),
	}
}

func (s *PokerServer) StartGame(_ context.Context, request *v1.StartGameRequest) (*v1.OperationResponse, error) {
	if _, ok := s.Players[request.GameId]; ok {
		return nil, errors.New("cannot create a game that already exists")
	}
	var newGame engine.Game
	if !request.Deterministic {
		newGame = engine.NewGame(int(request.SmallBlind), int(request.BigBlind))
	} else {
		newGame = engine.NewDeterministicGame(int(request.SmallBlind), int(request.BigBlind))
	}
	s.Games[request.GameId] = &newGame
	return &v1.OperationResponse{Successful: true, Message: "Successfully created game"}, nil
}

func (s *PokerServer) AddPlayer(_ context.Context, request *v1.AddPlayerRequest) (*v1.OperationResponse, error) {
	if _, ok := s.Players[request.PlayerId]; ok {
		return nil, errors.New("cannot add player that is already in game")
	}
	newPlayer := engine.Player{Name: request.Name, Stack: int(request.Stack)}
	s.Players[request.PlayerId] = &newPlayer
	return &v1.OperationResponse{Successful: true, Message: fmt.Sprintf("Successfully added '%s'", newPlayer.Name)}, nil
}

func (s *PokerServer) SitPlayer(_ context.Context, request *v1.SitPlayerRequest) (*v1.OperationResponse, error) {
	player, ok := s.Players[request.PlayerId]
	if !ok {
		return nil, errors.New(fmt.Sprintf("player with Player ID '%d' not found", request.PlayerId))
	}
	game, ok := s.Games[request.GameId]
	if !ok {
		return nil, errors.New(fmt.Sprintf("game with Game ID '%d' not found", request.GameId))
	}
	err := game.SitPlayer(player, int(request.SeatNumber))
	if err != nil {
		return nil, err
	}
	return &v1.OperationResponse{Successful: true, Message: fmt.Sprintf("Succesfully sat '%s' in seat %d", player.Name, player.SeatNumber)}, nil
}

func (s *PokerServer) StandPlayer(_ context.Context, request *v1.StandPlayerRequest) (*v1.OperationResponse, error) {
	player, ok := s.Players[request.PlayerId]
	if !ok {
		return nil, errors.New(fmt.Sprintf("player with Player ID '%d' not found", request.PlayerId))
	}
	game, ok := s.Games[request.GameId]
	if !ok {
		return nil, errors.New(fmt.Sprintf("game with Game ID '%d' not found", request.GameId))
	}
	err := game.StandPlayer(player, int(request.SeatNumber))
	if err != nil {
		return nil, err
	}
	return &v1.OperationResponse{Successful: true, Message: fmt.Sprintf("Succesfully stood '%s' from seat %d", player.Name, player.SeatNumber)}, nil
}

func (s *PokerServer) DealHand(_ context.Context, request *v1.DealHandRequest) (*v1.OperationResponse, error) {
	game, ok := s.Games[request.GameId]
	if !ok {
		return nil, errors.New(fmt.Sprintf("game with Game ID '%d' not found", request.GameId))
	}
	err := game.DealHand()
	if err != nil {
		return nil, err
	}
	return &v1.OperationResponse{Successful: true, Message: "Dealing hand"}, nil
}

func (s *PokerServer) TakeAction(_ context.Context, request *v1.TakeActionRequest) (*v1.OperationResponse, error) {
	player, ok := s.Players[request.PlayerId]
	if !ok {
		return nil, errors.New(fmt.Sprintf("player with Player ID '%d' not found", request.PlayerId))
	}
	game, ok := s.Games[request.GameId]
	if !ok {
		return nil, errors.New(fmt.Sprintf("game with Game ID '%d' not found", request.GameId))
	}
	err := game.TakeAction(player, engine.Action{
		ActionType: engine.ActionType(request.Action.ActionType),
		Value:      int(request.Action.Value),
	})
	if err != nil {
		return nil, err
	}
	return &v1.OperationResponse{Successful: true, Message: "Took action from player"}, nil
}

func (s *PokerServer) GetPlayerState(_ context.Context, request *v1.GetPlayerStateRequest) (*v1.PlayerState, error) {
	player, ok := s.Players[request.PlayerId]
	if !ok {
		return nil, errors.New(fmt.Sprintf("player with Player ID '%d' not found", request.PlayerId))
	}
	game, ok := s.Games[request.GameId]
	if !ok {
		return nil, errors.New(fmt.Sprintf("game with Game ID '%d' not found", request.GameId))
	}
	return game.GetPlayerState(player), nil
}

func StartDevServer() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 1234))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pokerServer := NewPokerServer()
	v1.RegisterPokerServiceServer(grpcServer, &pokerServer)
	reflection.Register(grpcServer)
	_ = grpcServer.Serve(lis)
}
