package engine

import "github.com/simster7/PocketRockets/go_backend/api"

type ActionType int

const (
	Check ActionType = iota
	Call
	Fold
	Bet
	Blind
)

type Action struct {
	ActionType ActionType
	Value      int
}

func (a *Action) GetMessage() *api.Action {
	return &api.Action{
		ActionType: int32(a.ActionType),
		Value:      int32(a.Value),
	}
}
