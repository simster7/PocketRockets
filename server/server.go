package server

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/simster7/PocketRockets/api/v1"
	"github.com/simster7/PocketRockets/engine"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type PokerServer struct {
	Games    map[int32]*engine.Game
	Personas map[int32]*engine.Persona
}

func NewPokerServer() PokerServer {
	return PokerServer{
		Games:    make(map[int32]*engine.Game),
		Personas: make(map[int32]*engine.Persona),
	}
}

func (s *PokerServer) StartGame(_ context.Context, request *v1.StartGameRequest) (*v1.OperationResponse, error) {
	if _, ok := s.Personas[request.GameId]; ok {
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

func (s *PokerServer) AddPersona(_ context.Context, request *v1.AddPersonaRequest) (*v1.OperationResponse, error) {
	if _, ok := s.Personas[request.PlayerId]; ok {
		return nil, errors.New("cannot add persona that already exists")
	}
	newPlayer := engine.Persona{Name: request.Name}
	s.Personas[request.PlayerId] = &newPlayer
	return &v1.OperationResponse{Successful: true, Message: fmt.Sprintf("Successfully added '%s'", newPlayer.Name)}, nil
}

func (s *PokerServer) SitPlayer(_ context.Context, request *v1.SitPlayerRequest) (*v1.OperationResponse, error) {
	persona, ok := s.Personas[request.PlayerId]
	if !ok {
		return nil, fmt.Errorf("persona with Persona ID '%d' not found", request.PlayerId)
	}
	game, ok := s.Games[request.GameId]
	if !ok {
		return nil, fmt.Errorf("game with Game ID '%d' not found", request.GameId)
	}
	err := game.SitPlayer(persona.Name, 100, int(request.SeatNumber))
	if err != nil {
		return nil, err
	}
	return &v1.OperationResponse{Successful: true, Message: fmt.Sprintf("Succesfully sat '%s' in seat %d", persona.Name, request.SeatNumber)}, nil
}

func (s *PokerServer) StandPlayer(_ context.Context, request *v1.StandPlayerRequest) (*v1.OperationResponse, error) {
	player, ok := s.Personas[request.PlayerId]
	if !ok {
		return nil, fmt.Errorf("player with Player ID '%d' not found", request.PlayerId)
	}
	game, ok := s.Games[request.GameId]
	if !ok {
		return nil, fmt.Errorf("game with Game ID '%d' not found", request.GameId)
	}
	err := game.StandPlayer(int(request.SeatNumber))
	if err != nil {
		return nil, err
	}
	return &v1.OperationResponse{Successful: true, Message: fmt.Sprintf("Succesfully stood '%s' from seat %d", player.Name, request.SeatNumber)}, nil
}

func (s *PokerServer) DealHand(_ context.Context, request *v1.DealHandRequest) (*v1.OperationResponse, error) {
	game, ok := s.Games[request.GameId]
	if !ok {
		return nil, fmt.Errorf("game with Game ID '%d' not found", request.GameId)
	}
	err := game.DealHand()
	if err != nil {
		return nil, err
	}
	return &v1.OperationResponse{Successful: true, Message: "Dealing hand"}, nil
}

func (s *PokerServer) TakeAction(_ context.Context, request *v1.TakeActionRequest) (*v1.OperationResponse, error) {
	_, ok := s.Personas[request.PlayerId]
	if !ok {
		return nil, fmt.Errorf("player with Player ID '%d' not found", request.PlayerId)
	}
	game, ok := s.Games[request.GameId]
	if !ok {
		return nil, fmt.Errorf("game with Game ID '%d' not found", request.GameId)
	}
	err := game.TakeAction(engine.Action{
		ActionType: engine.ActionType(request.Action.ActionType),
		Value:      int(request.Action.Value),
	})
	if err != nil {
		return nil, err
	}
	return &v1.OperationResponse{Successful: true, Message: "Took action from player"}, nil
}

func (s *PokerServer) GetPlayerState(_ context.Context, request *v1.GetPlayerStateRequest) (*v1.PlayerState, error) {
	_, ok := s.Personas[request.PlayerId]
	if !ok {
		return nil, fmt.Errorf("player with Player ID '%d' not found", request.PlayerId)
	}
	_, ok = s.Games[request.GameId]
	if !ok {
		return nil, fmt.Errorf("game with Game ID '%d' not found", request.GameId)
	}
	return nil, nil
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
