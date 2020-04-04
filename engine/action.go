package engine

import "github.com/simster7/PocketRockets/backend/api/v1"

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

func (a *Action) GetMessage() *v1.Action {
	return &v1.Action{
		ActionType: int32(a.ActionType),
		Value:      int32(a.Value),
	}
}
