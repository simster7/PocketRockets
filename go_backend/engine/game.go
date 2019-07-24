package engine

import "errors"

type Seat struct {
	Index int
	Occupied bool
	Player Player
}

type Game struct {
	Seats []Seat
	ButtonPosition int
	SmallBlind int
	BigBlind int
	GameState GameState
	IsHandActive bool
}

func NewGame(smallBlind, bigBlind int) Game {
	return Game{
		Seats: emptyTable(),
		ButtonPosition: 0,
		SmallBlind: smallBlind,
		BigBlind: bigBlind,
		//GameState: {}, // TODO
		IsHandActive: false,
	}
}

func (g *Game) SitPlayer(player Player, seatNumber int) error {
	if seatNumber < 0 || seatNumber >= 9 {
		return errors.New("invalid seat number")
	}
	if g.Seats[seatNumber].Occupied {
		return errors.New("cannot sit player on an occupied seat")
	}
	g.Seats[seatNumber] = Seat{
		Index: seatNumber,
		Occupied: true,
		Player: player,
	}
	return nil
}

func (g *Game) StandPlayer(player Player, seatNumber int) error {
	if seatNumber < 0 || seatNumber >= 9 {
		return errors.New("invalid seat number")
	}
	if !g.Seats[seatNumber].Occupied {
		return errors.New("seat is already empty")
	}
	// TODO Fix this
	if g.Seats[seatNumber].Player != player {
		return errors.New("incorrect player/seat number combination")
	}
	g.Seats[seatNumber] = Seat{
		Index: seatNumber,
		Occupied: false,
		Player: nil,
	}
	return nil
}

func emptyTable() []Seat {
	var seats []Seat
	for i := 0; i < 9; i++ {
		seats = append(seats, Seat{
			Index: i,
			Occupied: false,
			Player: nil,
		})
	}
	return seats
}