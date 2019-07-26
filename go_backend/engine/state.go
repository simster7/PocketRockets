package engine

import "github.com/simster7/PocketRockets/go_backend/engine/evaluator"

type Round int

const (
	PreFlop Round = iota
	Flop
	Turn
	River
)

const (
	Showdown string = "Showdown"
	Folds    string = "Folds"
)

type ActionConsequence struct {
	ValidAction bool
	Message     string

	Seat       Seat
	PlayerFold bool
	PlayerBet  int

	// Ends hand
	EndsHand      bool
	Payoffs       map[Seat]int
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
		FoldVector:     getInitialFoldVector(&seats),
		BetVector:      getZeroBetVector(),
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

func (gs *GameState) TakeAction(action Action) ActionConsequence {
	switch action.ActionType {
	case fold:
		gs.FoldVector[gs.ActingPlayer] = true
	case check:
		if gs.getLeadAction().ActionType == bet {
			return ActionConsequence{ValidAction: false, Message: "Illegal game state: player can't check when there is a bet"}
		}
	case call:
		amountToCall := gs.BetVector[gs.LeadingPlayer] - gs.BetVector[gs.ActingPlayer]
		// TODO: Replace this with all-in logic
		if gs.Seats[gs.ActingPlayer].Player.Stack < amountToCall {
			return ActionConsequence{ValidAction: false, Message: "Illegal game state: player doesn't have enough chips to call"}
		}

	}
	if action.ActionType == fold {
	}
	return ActionConsequence{ValidAction: false}
}

func (gs *GameState) moveActingPlayer() {
	gs.ActingPlayer = (gs.ActingPlayer + 1) % 9
	for gs.FoldVector[gs.ActingPlayer] && !gs.isRoundOver() {
		gs.ActingPlayer = (gs.ActingPlayer + 1) % 9
	}

	if gs.isRoundOver() {
		gs.ActingPlayer = gs.getNActivePlayerIndexFromIndex(gs.ButtonPosition, 1)
		gs.LeadingPlayer = gs.getNActivePlayerIndexFromIndex(gs.ButtonPosition, 1)
		gs.Pots[len(gs.Pots)-1] += getSum(gs.BetVector)
		gs.BetVector = getZeroBetVector()
		gs.Round++
	}

	if gs.isHandOver() {
		gs.IsHandActive = false
	}
}

func (gs *GameState) processEndGame(consequence *ActionConsequence) {
	if onePlayerStanding, player := gs.isOnePlayerStanding(); onePlayerStanding {
		// TODO: Process all-in logic here
		consequence.EndsHand = true
		consequence.WinCondition = Folds
		consequence.Payoffs[gs.Seats[player]] = gs.Pots[len(gs.Pots)-1]
	} else {
		
	}
}

func (gs *GameState) checkEndGame(consequence *ActionConsequence) {
	if !gs.isHandOver() {
		return
	}

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

func (gs *GameState) isRoundOver() bool {
	// TODO hard-code option for big blind when fold-around
	return gs.ActingPlayer == gs.LeadingPlayer
}

func (gs *GameState) isHandOver() bool {
	onePlayerStanding, _ := gs.isOnePlayerStanding()
	return (gs.isRoundOver() && gs.Round == River) || onePlayerStanding
}

func (gs *GameState) isOnePlayerStanding() (bool, int) {
	playersInHand := 0
	player := -1
	for i, folded := range gs.FoldVector {
		if !folded {
			playersInHand++
			player = i
		}
	}
	return playersInHand == 1, player
}

func (gs *GameState) getLeadAction() Action {
	return gs.Seats[gs.LeadingPlayer].Player.LastAction;
}

func getInitialFoldVector(seats *[9]Seat) [9]bool {
	var foldVector [9]bool
	for i, seat := range seats {
		foldVector[i] = !seat.Occupied || seat.Player.SittingOut
	}
	return foldVector
}

func getZeroBetVector() [9]int {
	return [9]int{0, 0, 0, 0, 0, 0, 0, 0, 0}
}

func getSum(a [9]int) int {
	count := 0
	for _, val := range a {
		count += val
	}
	return count
}
