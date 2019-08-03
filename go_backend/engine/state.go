package engine

import (
	"container/heap"
	"log"
)

type Round int

const (
	PreFlop Round = iota
	Flop
	Turn
	River
	HandEnd
)

const (
	Showdown string = "Showdown"
	Folds    string = "Folds"
)

var AllPlayers = []int{0, 1, 2, 3, 4, 5, 6, 7, 8}

type ActionConsequence struct {
	ValidAction bool
	Message     string

	PlayerIndex int
	PlayerFold  bool
	PlayerBet   int
	IsAllIn     bool

	// Refund over-bet money that can't be matched
	RefundsMoney      bool
	RefundPlayerIndex int
	RefundAmount      int

	// Ends hand
	EndsHand      bool
	Payoffs       map[Seat]int
	PotRemainder  int
	WinCondition  string
	ShowdownHands []HandForEvaluation
}

type BetVectorNode struct {
	Amount  int
	IsAllIn bool
}

type GameState struct {
	Seats          [9]Seat
	ButtonPosition int
	FoldVector     [9]bool
	BetVector      [9]BetVectorNode
	Pots           []int
	PotContenders  [][]int
	Deck           [52]Card
	Round          Round
	ActingPlayer   int
	LeadingPlayer  int
	IsHandActive   bool

	currentAction Action
}

func GetNewHandGameState(seats [9]Seat, buttonPosition, bigBlind, smallBlind int, deck [52]Card) (GameState, []ActionConsequence) {
	newState := GameState{
		Seats:          seats,
		ButtonPosition: buttonPosition,
		FoldVector:     getInitialFoldVector(&seats),
		BetVector:      getZeroBetVector(),
		Pots:           []int{0},
		PotContenders:  [][]int{AllPlayers},
		Deck:           deck,
		Round:          PreFlop,
	}

	smallBlindIndex := newState.getNActivePlayerIndexFromIndex(buttonPosition, 1)
	bigBlindIndex := newState.getNActivePlayerIndexFromIndex(buttonPosition, 2)
	utgIndex := newState.getNActivePlayerIndexFromIndex(buttonPosition, 3)

	newState.BetVector[bigBlindIndex].Amount = bigBlind
	newState.BetVector[smallBlindIndex].Amount = smallBlind

	newState.ActingPlayer = utgIndex
	newState.LeadingPlayer = bigBlindIndex

	newState.IsHandActive = true

	return newState, []ActionConsequence{
		{
			EndsHand:    false,
			ValidAction: true,
			PlayerIndex: bigBlindIndex,
			PlayerBet:   bigBlind,
		},
		{
			EndsHand:    false,
			ValidAction: true,
			PlayerIndex: smallBlindIndex,
			PlayerBet:   smallBlind,
		},
	}
}

func (gs *GameState) TakeAction(action Action) ActionConsequence {
	gs.currentAction = action
	var actionConsequence ActionConsequence
	switch action.ActionType {
	case Fold:
		gs.FoldVector[gs.ActingPlayer] = true
		actionConsequence = ActionConsequence{
			ValidAction: true,
			PlayerIndex: gs.ActingPlayer,
			PlayerFold:  true,
			IsAllIn:     false,
			PlayerBet:   0,
		}
	case Check:
		if gs.getLeadAction().ActionType == Bet {
			return ActionConsequence{
				ValidAction: false,
				Message:     "Illegal game state: player can't check when there is a bet",
			}
		}
		actionConsequence = ActionConsequence{
			ValidAction: true,
			PlayerIndex: gs.ActingPlayer,
			PlayerFold:  false,
			IsAllIn:     false,
			PlayerBet:   0,
		}
	case Call:
		if gs.getLeadAction().ActionType == Check {
			return ActionConsequence{
				ValidAction: false,
				Message:     "Illegal game state: player can't call when there is nothing to call",
			}
		}
		isAllIn := false
		amountToCall := gs.BetVector[gs.LeadingPlayer].Amount - gs.BetVector[gs.ActingPlayer].Amount
		if gs.Seats[gs.ActingPlayer].Player.Stack < amountToCall {
			amountToCall = gs.Seats[gs.ActingPlayer].Player.Stack
			gs.BetVector[gs.ActingPlayer].IsAllIn = true
			isAllIn = true
		}
		gs.BetVector[gs.ActingPlayer].Amount += amountToCall
		actionConsequence = ActionConsequence{
			ValidAction: true,
			PlayerIndex: gs.ActingPlayer,
			PlayerFold:  false,
			IsAllIn:     isAllIn,
			PlayerBet:   amountToCall,
		}
	case Bet:
		leadAction := gs.getLeadAction()
		toCall := 0
		isAllIn := false
		if leadAction.ActionType == Bet || leadAction.ActionType == Blind && leadAction.Value-gs.BetVector[gs.ActingPlayer].Amount > 0 {
			toCall = leadAction.Value - gs.BetVector[gs.ActingPlayer].Amount
		}
		if gs.Seats[gs.ActingPlayer].Player.Stack < action.Value {
			return ActionConsequence{
				ValidAction: false,
				Message:     "Illegal game state: player doesn't have enough chips to bet, should call all-in",
			}
		}
		if action.Value == (gs.Seats[gs.ActingPlayer].Player.Stack - toCall) {
			gs.BetVector[gs.ActingPlayer].IsAllIn = true
			isAllIn = true
		}
		callAndBet := toCall + action.Value
		gs.BetVector[gs.ActingPlayer].Amount += callAndBet
		gs.LeadingPlayer = gs.ActingPlayer
		actionConsequence = ActionConsequence{
			ValidAction: true,
			PlayerIndex: gs.ActingPlayer,
			PlayerFold:  false,
			IsAllIn:     isAllIn,
			PlayerBet:   callAndBet,
		}
	}
	endsHand := gs.moveActingPlayer(&actionConsequence)
	if endsHand {
		gs.processEndGame(&actionConsequence)
	}
	return actionConsequence
}

func (gs *GameState) moveActingPlayer(consequence *ActionConsequence) bool {
	gs.ActingPlayer = (gs.ActingPlayer + 1) % 9
	for (gs.FoldVector[gs.ActingPlayer] || gs.Seats[gs.ActingPlayer].Player.IsAllIn) && !gs.isRoundOver() {
		gs.ActingPlayer = (gs.ActingPlayer + 1) % 9
	}

	if gs.isRoundOver() {
		gs.ActingPlayer = gs.getNActivePlayerIndexFromIndex(gs.ButtonPosition, 1)
		gs.LeadingPlayer = gs.getNActivePlayerIndexFromIndex(gs.ButtonPosition, 1)
		gs.Seats[gs.LeadingPlayer].Player.LastAction = Action{ActionType: Check}
		processPots(&gs.BetVector, &gs.PotContenders, &gs.Pots, consequence)
		gs.Round++
	}

	if gs.isHandOver() {
		gs.IsHandActive = false
		return true
	}
	return false
}

func processPots(betVector *[9]BetVectorNode, potContenders *[][]int, pots *[]int, consequence *ActionConsequence) {
	var processPotsPQ ProcessPotsPQ
	totalRoundPot := 0
	highestBetAmountPlayerIndex := -1
	highestBetAmount := 0
	secondHighestBetAmount := 0
	isHighestBetCalled := false
	for i, node := range betVector {
		totalRoundPot += node.Amount
		if node.IsAllIn {
			processPotsPQ = append(processPotsPQ, &ProcessPotsPQItem{
				playerIndex: i,
				allInAmount: node.Amount,
				index:       len(processPotsPQ),
			})
		}
		if node.Amount == highestBetAmount {
			isHighestBetCalled = true
		} else if node.Amount > highestBetAmount {
			isHighestBetCalled = false
			secondHighestBetAmount = highestBetAmount
			highestBetAmount = node.Amount
			highestBetAmountPlayerIndex = i
		} else if node.Amount > secondHighestBetAmount {
			secondHighestBetAmount = node.Amount
		}
	}

	// One player's bet wasn't fully called (because of an all-in). Refund over-bet money.
	if highestBetAmount >= 0 && !isHighestBetCalled {
		consequence.RefundsMoney = true
		consequence.RefundPlayerIndex = highestBetAmountPlayerIndex
		amountToRefund := highestBetAmount - secondHighestBetAmount
		consequence.RefundAmount = amountToRefund
		totalRoundPot -= amountToRefund
	}

	heap.Init(&processPotsPQ)

	allInAlreadyContributed := 0
	for len(processPotsPQ) > 0 {
		var allInsAtCurrentAmount []*ProcessPotsPQItem
		allInsAtCurrentAmount = append(allInsAtCurrentAmount, heap.Pop(&processPotsPQ).(*ProcessPotsPQItem))
		currentAllInValue := allInsAtCurrentAmount[0].allInAmount - allInAlreadyContributed

		for len(processPotsPQ) > 0 && processPotsPQ[0].allInAmount == currentAllInValue {
			allInsAtCurrentAmount = append(allInsAtCurrentAmount, heap.Pop(&processPotsPQ).(*ProcessPotsPQItem))
		}

		currentContendersPot := 0
		for i, node := range betVector {
			if node.Amount > 0 {
				amountAddedToThisContention := min(node.Amount, currentAllInValue)
				currentContendersPot += amountAddedToThisContention
				betVector[i].Amount -= amountAddedToThisContention
				totalRoundPot -= amountAddedToThisContention
			}
		}

		currentContenders := (*potContenders)[len(*potContenders)-1]
		newContenders := filterInt(currentContenders, func(i int) bool {
			return !containsIntInProcessPotsPQSlice(allInsAtCurrentAmount, i)
		})
		(*pots)[len(*pots)-1] += currentContendersPot
		*pots = append(*pots, 0)
		*potContenders = append(*potContenders, newContenders)
		allInAlreadyContributed += currentAllInValue
	}
	(*pots)[len(*pots)-1] += totalRoundPot
	*betVector = getZeroBetVector()
}

func (gs *GameState) processEndGame(consequence *ActionConsequence) {
	consequence.Payoffs = make(map[Seat]int)

	if len(gs.Pots) != len(gs.PotContenders) {
		log.Fatal("bug: Pots and PotContenders are out of sync")
	}

	// Go through the side pots from latest to earliest
	for potIndex := len(gs.Pots) - 1; potIndex >= 0; potIndex-- {
		if onePlayerStanding, player := gs.isOnePlayerStanding(gs.PotContenders[potIndex]); onePlayerStanding {
			consequence.EndsHand = true
			consequence.WinCondition = Folds
			consequence.Payoffs[gs.Seats[player]] += gs.Pots[potIndex]
		} else {
			var showdownHands []HandForEvaluation
			communityCards := gs.getCommunityCards()
			for i, seat := range gs.Seats {
				if seat.Occupied && !seat.Player.SittingOut && gs.FoldVector[i] == false && containsIntInIntSlice(gs.PotContenders[potIndex], i) {
					showdownHands = append(showdownHands, HandForEvaluation{
						Hand:        append(gs.getPlayerCards(i), communityCards...),
						PlayerIndex: i,
					})
				}
			}
			rankedHands := EvaluateHands(showdownHands)
			consequence.EndsHand = true
			consequence.WinCondition = Showdown
			consequence.ShowdownHands = rankedHands

			numberOfWinners := 1
			for i := 1; i < len(rankedHands) && CompareStrengths(rankedHands[0].HandStrength, rankedHands[i].HandStrength) == 0; i++ {
				numberOfWinners++
			}

			for i := 0; i < numberOfWinners; i++ {
				consequence.Payoffs[gs.Seats[rankedHands[i].PlayerIndex]] += gs.Pots[potIndex] / numberOfWinners
			}
			consequence.PotRemainder += gs.Pots[potIndex] % numberOfWinners
		}
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
	// If limps around to big blind, give them option
	if gs.Round == PreFlop {
		bigBlindIndex := gs.getNActivePlayerIndexFromIndex(gs.ButtonPosition, 2)
		// Limps to big blind, give option
		if gs.ActingPlayer == bigBlindIndex && gs.LeadingPlayer == bigBlindIndex && gs.Seats[bigBlindIndex].Player.LastAction.ActionType == Blind {
			return false
		}
		// Big blind checked, end round
		utgPlayerIndex := gs.getNActivePlayerIndexFromIndex(gs.ButtonPosition, 3)
		if gs.ActingPlayer == utgPlayerIndex && gs.LeadingPlayer == bigBlindIndex && gs.currentAction.ActionType == Check {
			return true
		}
	}
	return gs.ActingPlayer == gs.LeadingPlayer
}

func (gs *GameState) isHandOver() bool {
	onePlayerStanding, _ := gs.isOnePlayerStanding(AllPlayers)
	return (gs.isRoundOver() && gs.Round == HandEnd) || onePlayerStanding
}

func (gs *GameState) isOnePlayerStanding(playersToConsider []int) (bool, int) {
	playersInHand := 0
	player := -1
	for i, folded := range gs.FoldVector {
		if !folded && containsIntInIntSlice(playersToConsider, i) {
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
	return []Card{gs.Deck[gs.getPlayerIndexInHand(playerIndex)], gs.Deck[gs.getPlayerIndexInHand(playerIndex)+gs.getNumberOfPlayersInHand()]}
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

// Small blind is 0
func (gs *GameState) getPlayerIndexInHand(seatIndex int) int {
	count := 0
	current := (gs.ButtonPosition + 1) % 9
	for i := 0; i < 9; i++ {
		if current == seatIndex {
			return count
		}
		if gs.Seats[current].Occupied {
			count++
		}
		current = (current + 1) % 9
	}
	log.Fatal("unreachable: getPlayerIndexInHand got a seatIndex that is not active: ", seatIndex)
	return 0
}

func getInitialFoldVector(seats *[9]Seat) [9]bool {
	var foldVector [9]bool
	for i, seat := range seats {
		foldVector[i] = !seat.Occupied || seat.Player.SittingOut
	}
	return foldVector
}

func getZeroBetVector() [9]BetVectorNode {
	var betVector [9]BetVectorNode
	for i := range betVector {
		betVector[i] = BetVectorNode{
			Amount:  0,
			IsAllIn: false,
		}
	}
	return betVector
}

func getSum(a [9]BetVectorNode) int {
	count := 0
	for _, node := range a {
		count += node.Amount
	}
	return count
}
