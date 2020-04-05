package engine

import (
	"container/heap"
	"log"
)

type State struct {
	Players        [9]*Player
	ButtonPosition int
	Pots           []int
	PotContenders  [][]int
	Deck           Deck
	Round          Round
	ActingPlayer   int
	LeadingPlayer  int
	IsHandActive   bool
	currentAction  Action
}

func GetNewHandState(players [9]*Player, buttonPosition, bigBlind, smallBlind int, deck Deck) (State, []ActionConsequence) {
	newState := State{
		Players:        players,
		ButtonPosition: buttonPosition,
		Pots:           []int{0},
		PotContenders:  [][]int{AllPlayers},
		Deck:           deck,
		Round:          RoundPreFlop,
	}

	smallBlindIndex := newState.getNActivePlayerIndexFromIndex(buttonPosition, 1)
	bigBlindIndex := newState.getNActivePlayerIndexFromIndex(buttonPosition, 2)
	utgIndex := newState.getNActivePlayerIndexFromIndex(buttonPosition, 3)

	newState.Players[bigBlindIndex].Bet = bigBlind
	newState.Players[smallBlindIndex].Bet = smallBlind

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

//func (gs *State) GetPlayerState(player *Player) *v1.PlayerState {
//	seatsMessage := make([]*v1.Seat, len(gs.Seats))
//	for i := 0; i < len(gs.Seats); i++ {
//		seatsMessage[i] = gs.Seats[i].GetMessage()
//	}
//
//	return &v1.PlayerState{
//		ButtonPosition: int32(gs.ButtonPosition),
//		BettingRound:   string(gs.Round),
//		Pots:           intSliceToInt32Slice(gs.Pots),
//		LeadPlayer:     int32(gs.LeadingPlayer),
//		ActingPlayer:   int32(gs.ActingPlayer),
//		Seats:          seatsMessage,
//		PlayerCards:    cardSliceToInt32Slice(gs.getPlayerCards(player.SeatNumber)),
//		CommunityCards: cardSliceToInt32Slice(gs.getCommunityCards()),
//		IsHandActive:   gs.IsHandActive,
//	}
//}

func (gs *State) TakeAction(action Action) ActionConsequence {
	gs.currentAction = action
	var actionConsequence ActionConsequence
	switch action.ActionType {
	case ActionTypeFold:
		gs.Players[gs.ActingPlayer].Folded = true
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
		amountToCall := gs.Players[gs.LeadingPlayer].Bet - gs.Players[gs.ActingPlayer].Bet
		if gs.Players[gs.ActingPlayer].Stack < amountToCall {
			amountToCall = gs.Players[gs.ActingPlayer].Stack
			gs.Players[gs.ActingPlayer].IsAllIn = true
			isAllIn = true
		}
		gs.Players[gs.ActingPlayer].Bet += amountToCall
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
		if leadAction.ActionType == ActionTypeBet || leadAction.ActionType == ActionTypeBlind && leadAction.Value-gs.Players[gs.ActingPlayer].Bet > 0 {
			toCall = leadAction.Value - gs.Players[gs.ActingPlayer].Bet
		}
		if gs.Players[gs.ActingPlayer].Stack < action.Value {
			return ActionConsequence{
				ValidAction: false,
				Message:     "Illegal game state: player doesn't have enough chips to bet, should call all-in",
			}
		}
		if action.Value == (gs.Players[gs.ActingPlayer].Stack - toCall) {
			gs.Players[gs.ActingPlayer].IsAllIn = true
			isAllIn = true
		}
		callAndBet := toCall + action.Value
		gs.Players[gs.ActingPlayer].Bet += callAndBet
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
	for (gs.Players[gs.ActingPlayer].Folded || gs.Players[gs.ActingPlayer].IsAllIn) && !gs.isRoundOver() {
		gs.ActingPlayer = (gs.ActingPlayer + 1) % 9
	}

	if gs.isRoundOver() {
		gs.ActingPlayer = gs.getNActivePlayerIndexFromIndex(gs.ButtonPosition, 1)
		gs.LeadingPlayer = gs.getNActivePlayerIndexFromIndex(gs.ButtonPosition, 1)
		gs.Players[gs.LeadingPlayer].LastAction = Action{ActionType: ActionTypeCheck}
		gs.processPots(consequence)
		gs.Round = gs.Round.GetNextRound()
	}

	if gs.isHandOver() {
		gs.IsHandActive = false
		return true
	}
	return false
}

func (gs *State) processPots(consequence *ActionConsequence) {
	var processPotsPQ ProcessPotsPQ
	totalRoundPot := 0
	highestBetAmountPlayerIndex := -1
	highestBetAmount := 0
	secondHighestBetAmount := 0
	isHighestBetCalled := false
	for i, player := range gs.Players {
		totalRoundPot += player.Bet
		if player.IsAllIn {
			processPotsPQ = append(processPotsPQ, &ProcessPotsPQItem{
				playerIndex: i,
				allInAmount: player.Bet,
				index:       len(processPotsPQ),
			})
		}
		if player.Bet == highestBetAmount {
			isHighestBetCalled = true
		} else if player.Bet > highestBetAmount {
			isHighestBetCalled = false
			secondHighestBetAmount = highestBetAmount
			highestBetAmount = player.Bet
			highestBetAmountPlayerIndex = i
		} else if player.Bet > secondHighestBetAmount {
			secondHighestBetAmount = player.Bet
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
		for _, player := range gs.Players {
			if player.Bet > 0 {
				amountAddedToThisContention := min(player.Bet, currentAllInValue)
				currentContendersPot += amountAddedToThisContention
				player.Bet -= amountAddedToThisContention
				totalRoundPot -= amountAddedToThisContention
			}
		}

		currentContenders := gs.PotContenders[len(gs.PotContenders)-1]
		newContenders := filterInt(currentContenders, func(i int) bool {
			return !containsIntInProcessPotsPQSlice(allInsAtCurrentAmount, i)
		})
		gs.Pots[len(gs.Pots)-1] += currentContendersPot
		gs.Pots = append(gs.Pots, 0)
		gs.PotContenders = append(gs.PotContenders, newContenders)
		allInAlreadyContributed += currentAllInValue
	}
	gs.Pots[len(gs.Pots)-1] += totalRoundPot
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
			consequence.Payoffs[gs.Players[player]] += gs.Pots[potIndex]
		} else {
			var showdownHands []HandForEvaluation
			communityCards := gs.getCommunityCards()
			for i, player := range gs.Players {
				if player != nil && !player.SittingOut && gs.Players[i].Folded == false && containsIntInIntSlice(gs.PotContenders[potIndex], i) {
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
				consequence.Payoffs[gs.Players[rankedHands[i].PlayerIndex]] += gs.Pots[potIndex] / numberOfWinners
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
		for gs.Players[index] != nil || gs.Players[index].Folded {
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
		if gs.ActingPlayer == bigBlindIndex && gs.LeadingPlayer == bigBlindIndex && gs.Players[bigBlindIndex].LastAction.ActionType == ActionTypeBlind {
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
	return gs.Players[gs.LeadingPlayer].LastAction
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
	for _, player := range gs.Players {
		if player != nil && !player.SittingOut {
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
		if gs.Players[current] != nil {
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
