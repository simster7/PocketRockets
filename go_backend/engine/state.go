package engine

import "github.com/simster7/PocketRockets/go_backend/engine/evaluator"

type Round int

const (
	PreFlop Round = iota
	Flop
	Turn
	River
)

type ActionConsequence struct {
	ValidAction bool

	Seat       Seat
	PlayerFold bool
	PlayerBet  int

	// Ends hand
	EndsHand      bool
	Payoffs       map[Player]int
	WinCondition  string
	ShowdownHands []evaluator.HandForEvaluation
}

type GameState struct {
	Seats          [9]Seat
	ButtonPosition int
	FoldVector     [9]bool
	BetVector      [9]int
	Pots           map[int]int
	PotContenders  map[int][]Player
	Deck           []Card
	Round          Round
	ActingPlayer   int
	LeadingPlayer  int
	IsHandActive   bool
}

func GetNewHandGameState(seats [9]Seat, buttonPosition, bigBlind, smallBlind int, deck []Card) (GameState, []ActionConsequence) {
	newState := GameState{
		Seats:          seats,
		ButtonPosition: buttonPosition,
		FoldVector:     [9]bool{false, false, false, false, false, false, false, false, false},
		BetVector:      [9]int{0, 0, 0, 0, 0, 0, 0, 0, 0},
		Pots:           map[int]int{0: 0},
		Deck:           deck,
		Round:          PreFlop,
	}

	smallBlindIndex := newState.getNActivePlayerIndexFromIndex(buttonPosition, 1)
	bigBlindIndex := newState.getNActivePlayerIndexFromIndex(buttonPosition, 2)
	utgIndex := newState.getNActivePlayerIndexFromIndex(buttonPosition, 3)

	newState.BetVector[bigBlindIndex] = bigBlind
	newState.BetVector[smallBlindIndex] = smallBlind
	newState.Seats[bigBlindIndex].Player.LastAction = Action{bet, bigBlind}
	newState.Seats[smallBlindIndex].Player.LastAction = Action{bet, smallBlind}

	newState.ActingPlayer = utgIndex
	newState.LeadingPlayer = bigBlindIndex

	newState.IsHandActive = true

	return newState, []ActionConsequence{
		{
			EndsHand:    false,
			ValidAction: true,
			Seat:        newState.Seats[bigBlindIndex],
			PlayerBet:   bigBlind,
		},
		{
			EndsHand:    false,
			ValidAction: true,
			Seat:        newState.Seats[smallBlindIndex],
			PlayerBet:   smallBlind,
		},
	}
}

func TakeAction(oldState GameState, action Action) (GameState, ActionConsequence) {
	return GameState{}, ActionConsequence{ValidAction: false}
}

// Returns index of player that is `n` active players to the right of `base`
func (gs *GameState) getNActivePlayerIndexFromIndex(base, n int) int {
	index := base
	count := 0
	for count != n {
		index = (index + 1) % 9
		for !gs.Seats[index].Occupied || gs.FoldVector[index] {
			index = (index + 1) % 9
		}
		count += 1
	}
	return index
}
