package engine

import (
	"errors"
	v1 "github.com/simster7/PocketRockets/api/v1"
	"log"
)

type Game struct {
	Seats          [9]Seat
	ButtonPosition int
	SmallBlind     int
	BigBlind       int
	GameState      State
	IsHandActive   bool
	Shuffler       func() Deck
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

	if player.SeatNumber != actionConsequence.PlayerIndex {
		log.Fatal("bug: unreachable: only acting player can have action consequence")
	}

	player.SetLastAction(action)
	player.SetIsAllIn(actionConsequence.IsAllIn)
	player.SetFolded(actionConsequence.PlayerFold)
	err := player.MakeBet(actionConsequence.PlayerBet)
	if err != nil {
		log.Fatal("bug: unreachable: player must have had enough to bet")
	}

	// Some actions may result in a player getting money refunded (such as over-betting an all-in that can't be fully
	// called)
	if actionConsequence.RefundsMoney {
		g.Seats[actionConsequence.RefundPlayerIndex].Player.Stack += actionConsequence.RefundAmount
	}

	if actionConsequence.EndsHand {
		g.IsHandActive = false
		for seat, amt := range actionConsequence.Payoffs {
			g.Seats[seat.Index].Player.ReceivePot(amt)
		}
	}
	return nil
}

func (g *Game) GetPlayerState(player *Player) *v1.PlayerState {
	return g.GameState.GetPlayerState(player)
}

func (g *Game) DealHand() error {

	if g.IsHandActive {
		return errors.New("cannot deal a hand while one is active")
	}

	if g.numberActivePlayers() < 2 {
		return errors.New("cannot deal a hand when only one player is active")
	}

	g.moveButton()

	deck := g.Shuffler()

	gameState, actionConsequences := GetNewHandState(g.Seats, g.ButtonPosition, g.BigBlind, g.SmallBlind, deck)

	g.GameState = gameState
	g.IsHandActive = true
	for _, action := range actionConsequences {
		err := g.Seats[action.PlayerIndex].Player.MakeBet(action.PlayerBet)
		if err != nil {
			log.Fatal("bug: unreachable: player must have had enough to bet")
		}
		g.Seats[action.PlayerIndex].Player.LastAction = Action{ActionType: ActionTypeBlind, Value: action.PlayerBet}
	}
	return nil
}

func (g *Game) moveButton() {
	g.ButtonPosition = (g.ButtonPosition + 1) % 9
	for !g.Seats[g.ButtonPosition].Occupied || g.Seats[g.ButtonPosition].Player.SittingOut {
		g.ButtonPosition = (g.ButtonPosition + 1) % 9
	}
}

func (g *Game) numberActivePlayers() int {
	count := 0
	for _, seat := range g.Seats {
		if seat.Occupied && !seat.Player.SittingOut {
			count++
		}
	}
	return count
}
