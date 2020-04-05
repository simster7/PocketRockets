package engine

import (
	"errors"
)

type Game struct {
	Seats          Seats
	ButtonPosition int
	SmallBlind     int
	BigBlind       int
	GameState      State
	IsHandActive   bool
	Shuffler       func() Deck
}

func NewGame(smallBlind, bigBlind int) Game {
	return Game{
		Seats:        emptyTable(),
		SmallBlind:   smallBlind,
		BigBlind:     bigBlind,
		IsHandActive: false,
		Shuffler:     getShuffledDeck,
	}
}

func NewDeterministicGame(smallBlind, bigBlind int) Game {
	return Game{
		Seats:        emptyTable(),
		SmallBlind:   smallBlind,
		BigBlind:     bigBlind,
		IsHandActive: false,
		Shuffler:     getDeck,
	}
}

func (g *Game) SitPlayer(name string, stack int, seatNumber int) error {
	if seatNumber < 0 || seatNumber >= 9 {
		return errors.New("invalid seat number")
	}
	if g.Seats[seatNumber] != nil {
		return errors.New("cannot sit player on an occupied seat")
	}
	g.Seats[seatNumber] = &Seat{
		Name:     name,
		Stack:    stack,
	}
	return nil
}

func (g *Game) StandPlayer(player *Player, seatNumber int) error {
	if seatNumber < 0 || seatNumber >= 9 {
		return errors.New("invalid seat number")
	}
	if g.Seats[seatNumber] == nil {
		return errors.New("seat is already empty")
	}
	g.Seats[seatNumber] = nil
	return nil
}

func (g *Game) TakeAction(player *Player, action Action) error {
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

	g.moveButton()
	gameState := GetNewHandState(g.getActivePlayers(), g.ButtonPosition, g.BigBlind, g.SmallBlind, g.Shuffler())
	g.GameState = gameState
	g.IsHandActive = true
	return nil
}

func (g *Game) moveButton() {
	g.ButtonPosition = (g.ButtonPosition + 1) % 9
	for g.Seats[g.ButtonPosition] == nil || g.Seats[g.ButtonPosition].SittingOut {
		g.ButtonPosition = (g.ButtonPosition + 1) % 9
	}
}

func (g *Game) numberActivePlayers() int {
	count := 0
	for _, seat := range g.Seats {
		if seat != nil && !seat.SittingOut {
			count++
		}
	}
	return count
}

func (g *Game) getActivePlayers() Players {
	players := new(Players)
	for i, seat := range g.Seats {
		if seat != nil && !seat.SittingOut {
			players[i] = &Player{Name: seat.Name, Stack: seat.Stack}
		}
	}
	return *players
}
