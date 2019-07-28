package engine

import (
	"errors"
	"log"
	"math/rand"
)

type Seat struct {
	Index    int
	Occupied bool
	Player   *Player
}

type Game struct {
	Seats          [9]Seat
	PlayerSeats    map[Player]Seat
	ButtonPosition int
	SmallBlind     int
	BigBlind       int
	GameState      GameState
	IsHandActive   bool
	Shuffler       func() [52]Card
}

func NewGame(smallBlind, bigBlind int) Game {
	return Game{
		Seats:          emptyTable(),
		ButtonPosition: 0,
		SmallBlind:     smallBlind,
		BigBlind:       bigBlind,
		IsHandActive:   false,
		Shuffler:       getShuffledDeck,
	}
}

func NewDeterministicGame(smallBlind, bigBlind int) Game {
	return Game{
		Seats:          emptyTable(),
		ButtonPosition: 0,
		SmallBlind:     smallBlind,
		BigBlind:       bigBlind,
		IsHandActive:   false,
		Shuffler:       getDeck,
	}
}

func (g *Game) SitPlayer(player *Player, seatNumber int) error {
	if seatNumber < 0 || seatNumber >= 9 {
		return errors.New("invalid seat number")
	}
	if g.Seats[seatNumber].Occupied {
		return errors.New("cannot sit player on an occupied seat")
	}
	g.Seats[seatNumber] = Seat{
		Index:    seatNumber,
		Occupied: true,
		Player:   player,
	}
	player.SeatNumber = seatNumber
	return nil
}

func (g *Game) StandPlayer(player *Player, seatNumber int) error {
	if seatNumber < 0 || seatNumber >= 9 {
		return errors.New("invalid seat number")
	}
	if !g.Seats[seatNumber].Occupied {
		return errors.New("seat is already empty")
	}
	if player.SeatNumber != seatNumber {
		return errors.New("incorrect player/seat number combination")
	}
	g.Seats[seatNumber] = Seat{
		Index:    seatNumber,
		Occupied: false,
	}
	player.SeatNumber = -1
	return nil
}

func (g *Game) TakeAction(player *Player, action Action) error {
	if !g.IsHandActive {
		return errors.New("cannot take action when hand is not active")
	}
	if player.SeatNumber != g.GameState.ActingPlayer {
		return errors.New("player is out of turn")
	}

	actionConsequence := g.GameState.TakeAction(action)
	if actionConsequence.ValidAction == false {
		return errors.New(actionConsequence.Message)
	}

	if player.SeatNumber != actionConsequence.Seat.Index {
		log.Fatal("bug: unreachable: only acting player can have action consequence")
	}

	player.SetLastAction(action)
	player.SetFolded(actionConsequence.PlayerFold)
	err := player.MakeBet(actionConsequence.PlayerBet)
	if err != nil {
		log.Fatal("bug: unreachable: player must have had enough to bet")
	}

	if actionConsequence.EndsHand {
		g.IsHandActive = false
		for seat, amt := range actionConsequence.Payoffs {
			g.Seats[seat.Index].Player.ReceivePot(amt)
		}
	}
	return nil
}

func (g *Game) DealHand() {
	g.moveButton()

	deck := g.Shuffler()

	gameState, actionConsequences := GetNewHandGameState(g.Seats, g.ButtonPosition, g.BigBlind, g.SmallBlind, deck)

	g.GameState = gameState
	g.IsHandActive = true
	for _, action := range actionConsequences {
		err := g.Seats[action.Seat.Index].Player.MakeBet(action.PlayerBet)
		if err != nil {
			log.Fatal("bug: unreachable: player must have had enough to bet")
		}
		g.Seats[action.Seat.Index].Player.LastAction = Action{ActionType: bet, Value: action.PlayerBet}
	}
}

func (g *Game) moveButton() {
	g.ButtonPosition = (g.ButtonPosition + 1) % 9
	for !g.Seats[g.ButtonPosition].Occupied || g.Seats[g.ButtonPosition].Player.SittingOut {
		g.ButtonPosition = (g.ButtonPosition + 1) % 9
	}
}

func getShuffledDeck() [52]Card {
	var deck [52]Card
	perm := rand.Perm(52)
	for i := 0; i < 52; i++ {
		deck[perm[i]] = Card{i}
	}
	return deck
}

func getDeck() [52]Card {
	var deck [52]Card
	for i := 0; i < 52; i++ {
		deck[i] = Card{i}
	}
	return deck
}

func emptyTable() [9]Seat {
	var seats [9]Seat
	for i := 0; i < 9; i++ {
		seats[i] = Seat{
			Index:    i,
			Occupied: false,
			Player:   nil,
		}
	}
	return seats
}
