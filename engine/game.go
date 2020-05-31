package engine

import (
	"errors"
)

type Game struct {
	SmallBlind   int
	BigBlind     int
	GameState    State
	IsHandActive bool
	Shuffler     Shuffler
}

func NewGame(smallBlind, bigBlind int) Game {
	return Game{
		SmallBlind:   smallBlind,
		BigBlind:     bigBlind,
		GameState:    NewHandState(*new(Players), 0),
		IsHandActive: false,
		Shuffler:     &StandardShuffler{},
	}
}

func NewDeterministicGame(smallBlind, bigBlind int) Game {
	return Game{
		SmallBlind:   smallBlind,
		BigBlind:     bigBlind,
		GameState:    NewHandState(*new(Players), 0),
		IsHandActive: false,
		Shuffler:     &DeterministicShuffler{},
	}
}

func (g *Game) SitPlayer(name string, stack, seat int) error {
	if player := g.GameState.Players[seat]; player != nil {
		return errors.New("cannot sit player; seat is occupied")
	}
	player := Player{Name: name, Stack: stack, Waiting: false}
	g.GameState.Players[seat] = &player

	return nil
}
func (g *Game) StandPlayer(seat int) error {
	if player := g.GameState.Players[seat]; player == nil {
		return errors.New("cannot stand player; seat is not occupied")
	}
	g.GameState.Players[seat] = nil

	return nil
}

func (g *Game) TakeAction(action Action) error {
	if !g.IsHandActive {
		return errors.New("cannot take action when hand is not active")
	}
	return g.GameState.TakeAction(action)
}

func (g *Game) DealHand() error {
	if g.IsHandActive {
		return errors.New("cannot deal a hand while one is active")
	}
	if g.numberActivePlayers() < 2 {
		return errors.New("cannot deal a hand when only one player is active")
	}

	g.GameState.DealHand(g.BigBlind, g.SmallBlind, g.Shuffler.Shuffle())
	g.IsHandActive = true
	return nil
}

func (g *Game) numberActivePlayers() int {
	count := 0
	for _, player := range g.GameState.Players {
		if player != nil && !player.SittingOut {
			count++
		}
	}
	return count
}
