package server

import (
	"context"
	"fmt"
	"github.com/simster7/PocketRockets/go_backend/api/v1"
	"github.com/simster7/PocketRockets/go_backend/engine"
)

type PokerServer struct {

}

func (s *PokerServer) GetPlayerState(context.Context, *v1.PlayerStateRequest) (*v1.PlayerState, error) {
	fmt.Println("Called")
	return &v1.PlayerState{BettingRound: int32(engine.PreFlop), ButtonPosition: int32(3)}, nil
}