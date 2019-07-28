package engine

type Round int

const (
	PreFlop Round = iota
	Flop
	Turn
	River
	PostRiver
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
	ShowdownHands []HandForEvaluation
}

type GameState struct {
	Seats          [9]Seat
	ButtonPosition int
	FoldVector     [9]bool
	BetVector      [9]int
	Pots           []int
	PotContenders  map[int][]Player
	Deck           [52]Card
	Round          Round
	ActingPlayer   int
	LeadingPlayer  int
	IsHandActive   bool
}

func GetNewHandGameState(seats [9]Seat, buttonPosition, bigBlind, smallBlind int, deck [52]Card) (GameState, []ActionConsequence) {
	newState := GameState{
		Seats:          seats,
		ButtonPosition: buttonPosition,
		FoldVector:     getInitialFoldVector(&seats),
		BetVector:      getZeroBetVector(),
		Pots:           []int{0},
		Deck:           deck,
		Round:          PreFlop,
	}

	smallBlindIndex := newState.getNActivePlayerIndexFromIndex(buttonPosition, 1)
	bigBlindIndex := newState.getNActivePlayerIndexFromIndex(buttonPosition, 2)
	utgIndex := newState.getNActivePlayerIndexFromIndex(buttonPosition, 3)

	newState.BetVector[bigBlindIndex] = bigBlind
	newState.BetVector[smallBlindIndex] = smallBlind

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
	var actionConsequence ActionConsequence
	switch action.ActionType {
	case fold:
		gs.FoldVector[gs.ActingPlayer] = true
		actionConsequence = ActionConsequence{
			ValidAction: true,
			Seat:        gs.Seats[gs.ActingPlayer],
			PlayerFold:  true,
			PlayerBet:   0,
		}
	case check:
		if gs.getLeadAction().ActionType == bet {
			return ActionConsequence{
				ValidAction: false,
				Message:     "Illegal game state: player can't check when there is a bet",
			}
		}
		actionConsequence = ActionConsequence{
			ValidAction: true,
			Seat:        gs.Seats[gs.ActingPlayer],
			PlayerFold:  false,
			PlayerBet:   0,
		}
	case call:
		if gs.getLeadAction().ActionType == check {
			return ActionConsequence{
				ValidAction: false,
				Message:     "Illegal game state: player can't call when there is nothing to call",
			}
		}
		amountToCall := gs.BetVector[gs.LeadingPlayer] - gs.BetVector[gs.ActingPlayer]
		// TODO: Replace this with all-in logic
		if gs.Seats[gs.ActingPlayer].Player.Stack < amountToCall {
			return ActionConsequence{
				ValidAction: false,
				Message:     "Illegal game state: player doesn't have enough chips to call",
			}
		}
		gs.BetVector[gs.ActingPlayer] += amountToCall
		actionConsequence = ActionConsequence{
			ValidAction: true,
			Seat:        gs.Seats[gs.ActingPlayer],
			PlayerFold:  false,
			PlayerBet:   amountToCall,
		}
	case bet:
		leadAction := gs.getLeadAction()
		toCall := 0
		if leadAction.ActionType == bet && leadAction.Value-gs.BetVector[gs.ActingPlayer] > 0 {
			toCall = leadAction.Value - gs.BetVector[gs.ActingPlayer]
		}
		// TODO: All-in logic here
		if gs.Seats[gs.ActingPlayer].Player.Stack < action.Value {
			return ActionConsequence{
				ValidAction: false,
				Message:     "Illegal game state: player doesn't have enough chips to make that bet",
			}
		}
		callAndBet := toCall + action.Value
		gs.BetVector[gs.ActingPlayer] += callAndBet
		gs.LeadingPlayer = gs.ActingPlayer
		actionConsequence = ActionConsequence{
			ValidAction: true,
			Seat:        gs.Seats[gs.ActingPlayer],
			PlayerFold:  false,
			PlayerBet:   callAndBet,
		}
	}
	endsHand := gs.moveActingPlayer()
	if endsHand {
		gs.processEndGame(&actionConsequence)
	}
	return actionConsequence
}

func (gs *GameState) moveActingPlayer() bool {
	gs.ActingPlayer = (gs.ActingPlayer + 1) % 9
	for gs.FoldVector[gs.ActingPlayer] && !gs.isRoundOver() {
		gs.ActingPlayer = (gs.ActingPlayer + 1) % 9
	}

	if gs.isRoundOver() {
		gs.ActingPlayer = gs.getNActivePlayerIndexFromIndex(gs.ButtonPosition, 1)
		gs.LeadingPlayer = gs.getNActivePlayerIndexFromIndex(gs.ButtonPosition, 1)
		gs.Seats[gs.LeadingPlayer].Player.LastAction = Action{ActionType: check}
		gs.Pots[len(gs.Pots)-1] += getSum(gs.BetVector)
		gs.BetVector = getZeroBetVector()
		gs.Round++
	}

	if gs.isHandOver() {
		gs.IsHandActive = false
		return true
	}
	return false
}

func (gs *GameState) processEndGame(consequence *ActionConsequence) {
	consequence.Payoffs = make(map[Seat]int)
	if onePlayerStanding, player := gs.isOnePlayerStanding(); onePlayerStanding {
		// TODO: Process all-in logic here
		consequence.EndsHand = true
		consequence.WinCondition = Folds
		consequence.Payoffs[gs.Seats[player]] = gs.Pots[len(gs.Pots)-1]
	} else {
		var showdownHands []HandForEvaluation
		communityCards := gs.getCommunityCards()
		for i, seat := range gs.Seats {
			if seat.Occupied && !seat.Player.SittingOut && gs.FoldVector[i] == false {
				showdownHands = append(showdownHands, HandForEvaluation{
					Hand:        append(gs.getPlayerCards(i), communityCards...),
					PlayerIndex: i,
				})
			}
		}
		rankedHands := EvaluateHands(showdownHands)
		// TODO: Add split-pot and all-in logic here
		consequence.EndsHand = true
		consequence.WinCondition = Showdown
		consequence.Payoffs[gs.Seats[rankedHands[0].PlayerIndex]] = gs.Pots[len(gs.Pots)-1]
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
	return (gs.isRoundOver() && gs.Round == PostRiver) || onePlayerStanding
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
	return gs.Seats[gs.LeadingPlayer].Player.LastAction
}

func (gs *GameState) getPlayerCards(playerIndex int) []Card {
	return []Card{gs.Deck[playerIndex], gs.Deck[playerIndex+gs.getNumberOfPlayersInHand()]}
}

func (gs *GameState) getCommunityCards() []Card {
	var communityCards []Card
	numPlayers := gs.getNumberOfPlayersInHand()
	if gs.Round >= Flop {
		communityCards = append(communityCards, gs.Deck[2*numPlayers+1:2*numPlayers+4]...)
	}
	if gs.Round >= Turn {
		communityCards = append(communityCards, gs.Deck[2*numPlayers+5])
	}
	if gs.Round >= River {
		communityCards = append(communityCards, gs.Deck[2*numPlayers+7])
	}
	return communityCards
}

func (gs *GameState) getNumberOfPlayersInHand() int {
	count := 0
	for _, seat := range gs.Seats {
		if seat.Occupied && !seat.Player.SittingOut {
			count++
		}
	}
	return count
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
