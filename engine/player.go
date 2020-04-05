package engine

import (
	"errors"
	"github.com/simster7/PocketRockets/api/v1"
)

type Player struct {
	Name       string
	Stack      int
	SeatNumber int
	Folded     bool
	IsAllIn    bool
	LastAction Action
	SittingOut bool
}

func NewPlayer(name string, stack int) Player {
	return Player{
		Name:       name,
		Stack:      stack,
		SeatNumber: -1,
		Folded:     false,
		LastAction: Action{},
		SittingOut: false,
	}
}

func (p *Player) MakeBet(betSize int) error {
	if p.Stack < betSize {
		return errors.New("player may not bet amount larger than their stack")
	}
	p.Stack -= betSize
	return nil
}

func (p *Player) ReceivePot(potSize int) {
	p.Stack += potSize
}

func (p *Player) SetLastAction(lastAction Action) {
	p.LastAction = lastAction
}

func (p *Player) SetFolded(folded bool) {
	p.Folded = folded
}

func (p *Player) SetIsAllIn(isAllIn bool) {
	p.IsAllIn = isAllIn
}

func (p *Player) GetMessage() *v1.Player {
	return &v1.Player{
		Name:       p.Name,
		Stack:      int32(p.Stack),
		SeatNumber: int32(p.SeatNumber),
		Folded:     p.Folded,
		IsAllIn:    p.IsAllIn,
		SittingOut: p.SittingOut,
		LastAction: p.LastAction.GetMessage(),
	}
}
