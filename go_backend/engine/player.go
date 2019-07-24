package engine

import "errors"

type Player struct {
	Name       string
	Stack      int
	SeatNumber int
	Folded     bool
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