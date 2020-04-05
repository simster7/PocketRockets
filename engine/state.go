package engine

import (
	"container/heap"
	v1 "github.com/simster7/PocketRockets/api/v1"
	"log"
)

type State struct {
	Seats          [9]Seat
	ButtonPosition int
	FoldVector     FoldVector
	BetVector      BetVector
	Pots           []int
	PotContenders  [][]int
	Deck           Deck
	Round          Round
	ActingPlayer   int
	LeadingPlayer  int
	IsHandActive   bool
	currentAction  Action
}

func GetNewHandState(seats [9]Seat, buttonPosition, bigBlind, smallBlind int, deck Deck) (State, []ActionConsequence) {
	newState := State{
		Seats:          seats,
		ButtonPosition: buttonPosition,
		FoldVector:     getInitialFoldVector(&seats),
		BetVector:      getZeroBetVector(),
		Pots:           []int{0},
		PotContenders:  [][]int{AllPlayers},
		Deck:           deck,
		Round:          RoundPreFlop,
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

func (gs *State) GetPlayerState(player *Player) *v1.PlayerState {
	seatsMessage := make([]*v1.Seat, len(gs.Seats))
	for i := 0; i < len(gs.Seats); i++ {
		seatsMessage[i] = gs.Seats[i].GetMessage()
	}

	return &v1.PlayerState{
		ButtonPosition: int32(gs.ButtonPosition),
		BettingRound:   string(gs.Round),
		Pots:           intSliceToInt32Slice(gs.Pots),
		LeadPlayer:     int32(gs.LeadingPlayer),
		ActingPlayer:   int32(gs.ActingPlayer),
		Seats:          seatsMessage,
		PlayerCards:    cardSliceToInt32Slice(gs.getPlayerCards(player.SeatNumber)),
		CommunityCards: cardSliceToInt32Slice(gs.getCommunityCards()),
		IsHandActive:   gs.IsHandActive,
	}
}

func (gs *State) TakeAction(action Action) ActionConsequence {
	gs.currentAction = action
	var actionConsequence ActionConsequence
	switch action.ActionType {
	case ActionTypeFold:
		gs.FoldVector[gs.ActingPlayer] = true
		actionConsequence = ActionConsequence{
			ValidAction: true,
			PlayerIndex: gs.ActingPlayer,
			PlayerFold:  true,
			IsAllIn:     false,
			PlayerBet:   0,
		}
	case ActionTypeCheck:
		if gs.getLeadAction().ActionType == ActionTypeBet {
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
	case ActionTypeCall:
		if gs.getLeadAction().ActionType == ActionTypeCheck {
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
	case ActionTypeBet:
		leadAction := gs.getLeadAction()
		toCall := 0
		isAllIn := false
		if leadAction.ActionType == ActionTypeBet || leadAction.ActionType == ActionTypeBlind && leadAction.Value-gs.BetVector[gs.ActingPlayer].Amount > 0 {
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

func (gs *State) moveActingPlayer(consequence *ActionConsequence) bool {
	gs.ActingPlayer = (gs.ActingPlayer + 1) % 9
	for (gs.FoldVector[gs.ActingPlayer] || gs.Seats[gs.ActingPlayer].Player.IsAllIn) && !gs.isRoundOver() {
		gs.ActingPlayer = (gs.ActingPlayer + 1) % 9
	}

	if gs.isRoundOver() {
		gs.ActingPlayer = gs.getNActivePlayerIndexFromIndex(gs.ButtonPosition, 1)
		gs.LeadingPlayer = gs.getNActivePlayerIndexFromIndex(gs.ButtonPosition, 1)
		gs.Seats[gs.LeadingPlayer].Player.LastAction = Action{ActionType: ActionTypeCheck}
		processPots(&gs.BetVector, &gs.PotContenders, &gs.Pots, consequence)
		gs.Round = gs.Round.GetNextRound()
	}

	if gs.isHandOver() {
		gs.IsHandActive = false
		return true
	}
	return false
}

func processPots(betVector *BetVector, potContenders *[][]int, pots *[]int, consequence *ActionConsequence) {
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

func (gs *State) processEndGame(consequence *ActionConsequence) {
	consequence.Payoffs = make(map[Seat]int)

	if len(gs.Pots) != len(gs.PotContenders) {
		log.Fatal("bug: Pots and PotContenders are out of sync")
	}

	// Go through the side pots from latest to earliest
	for potIndex := len(gs.Pots) - 1; potIndex >= 0; potIndex-- {
		if onePlayerStanding, player := gs.isOnePlayerStanding(gs.PotContenders[potIndex]); onePlayerStanding {
			consequence.EndsHand = true
			consequence.WinCondition = WinConditionFolds
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
			consequence.WinCondition = WinConditionShowdown
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
func (gs *State) getNActivePlayerIndexFromIndex(base, n int) int {
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

func (gs *State) isRoundOver() bool {
	// If limps around to big blind, give them option
	if gs.Round == RoundPreFlop {
		bigBlindIndex := gs.getNActivePlayerIndexFromIndex(gs.ButtonPosition, 2)
		// Limps to big blind, give option
		if gs.ActingPlayer == bigBlindIndex && gs.LeadingPlayer == bigBlindIndex && gs.Seats[bigBlindIndex].Player.LastAction.ActionType == ActionTypeBlind {
			return false
		}
		// Big blind checked, end round
		utgPlayerIndex := gs.getNActivePlayerIndexFromIndex(gs.ButtonPosition, 3)
		if gs.ActingPlayer == utgPlayerIndex && gs.LeadingPlayer == bigBlindIndex && gs.currentAction.ActionType == ActionTypeCheck {
			return true
		}
	}
	return gs.ActingPlayer == gs.LeadingPlayer
}

func (gs *State) isHandOver() bool {
	onePlayerStanding, _ := gs.isOnePlayerStanding(AllPlayers)
	return (gs.isRoundOver() && gs.Round == RoundHandEnd) || onePlayerStanding
}

func (gs *State) isOnePlayerStanding(playersToConsider []int) (bool, int) {
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

func (gs *State) getLeadAction() Action {
	return gs.Seats[gs.LeadingPlayer].Player.LastAction
}

func (gs *State) getPlayerCards(playerIndex int) []Card {
	return []Card{gs.Deck[gs.getPlayerIndexInHand(playerIndex)], gs.Deck[gs.getPlayerIndexInHand(playerIndex)+gs.getNumberOfPlayersInHand()]}
}

func (gs *State) getCommunityCards() []Card {
	var communityCards []Card
	numPlayers := gs.getNumberOfPlayersInHand()
	if gs.Round.IsAtOrAfter(RoundFlop) {
		communityCards = append(communityCards, gs.Deck[2*numPlayers+1:2*numPlayers+4]...)
	}
	if gs.Round.IsAtOrAfter(RoundTurn) {
		communityCards = append(communityCards, gs.Deck[2*numPlayers+5])
	}
	if gs.Round.IsAtOrAfter(RoundRiver) {
		communityCards = append(communityCards, gs.Deck[2*numPlayers+7])
	}
	return communityCards
}

func (gs *State) getNumberOfPlayersInHand() int {
	count := 0
	for _, seat := range gs.Seats {
		if seat.Occupied && !seat.Player.SittingOut {
			count++
		}
	}
	return count
}

// Small blind is 0
func (gs *State) getPlayerIndexInHand(seatIndex int) int {
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

func getZeroBetVector() BetVector {
	var betVector BetVector
	for i := range betVector {
		betVector[i] = BetVectorNode{
			Amount:  0,
			IsAllIn: false,
		}
	}
	return betVector
}

func getSum(a BetVector) int {
	count := 0
	for _, node := range a {
		count += node.Amount
	}
	return count
}
