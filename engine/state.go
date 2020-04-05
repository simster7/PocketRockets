package engine

import (
	"container/heap"
	"errors"
	"log"
)

type State struct {
	Players        Players
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

func GetNewHandState(players Players, buttonPosition, bigBlind, smallBlind int, deck Deck) State {
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

	return newState
}

func (gs *State) TakeAction(action Action) error {
	gs.currentAction = action
	switch action.ActionType {
	case ActionTypeFold:
		gs.Players[gs.ActingPlayer].Folded = true
	case ActionTypeCheck:
		if gs.getLeadAction().ActionType == ActionTypeBet {
			return errors.New("invalid action: player can't check when there is a bet")
		}
	case ActionTypeCall:
		if gs.getLeadAction().ActionType == ActionTypeCheck {
			return errors.New("invalid action: player can't call when there is nothing to call")
		}
		amountToCall := gs.Players[gs.LeadingPlayer].Bet - gs.Players[gs.ActingPlayer].Bet
		if gs.Players[gs.ActingPlayer].Stack < amountToCall {
			amountToCall = gs.Players[gs.ActingPlayer].Stack
			gs.Players[gs.ActingPlayer].IsAllIn = true
		}
		gs.Players[gs.ActingPlayer].Bet += amountToCall
	case ActionTypeBet:
		leadAction := gs.getLeadAction()
		toCall := 0
		if leadAction.ActionType == ActionTypeBet || leadAction.ActionType == ActionTypeBlind && leadAction.Value-gs.Players[gs.ActingPlayer].Bet > 0 {
			toCall = leadAction.Value - gs.Players[gs.ActingPlayer].Bet
		}
		if gs.Players[gs.ActingPlayer].Stack < action.Value {
			return errors.New("invalid action: player doesn't have enough chips to bet, should call all-in")
		}
		if action.Value == (gs.Players[gs.ActingPlayer].Stack - toCall) {
			gs.Players[gs.ActingPlayer].IsAllIn = true
		}
		callAndBet := toCall + action.Value
		gs.Players[gs.ActingPlayer].Bet += callAndBet
		gs.LeadingPlayer = gs.ActingPlayer
	}
	gs.Players[gs.ActingPlayer].LastAction = action
	endsHand := gs.moveActingPlayer()
	if endsHand {
		gs.processEndGame()
	}
	return nil
}

func (gs *State) moveActingPlayer() bool {
	gs.ActingPlayer = (gs.ActingPlayer + 1) % 9
	for (gs.Players[gs.ActingPlayer] == nil || gs.Players[gs.ActingPlayer].Folded || gs.Players[gs.ActingPlayer].IsAllIn) && !gs.isRoundOver() {
		gs.ActingPlayer = (gs.ActingPlayer + 1) % 9
	}

	if gs.isRoundOver() {
		gs.ActingPlayer = gs.getNActivePlayerIndexFromIndex(gs.ButtonPosition, 1)
		gs.LeadingPlayer = gs.getNActivePlayerIndexFromIndex(gs.ButtonPosition, 1)
		gs.Players[gs.LeadingPlayer].LastAction = Action{ActionType: ActionTypeCheck}
		gs.processPots()
		gs.Round = gs.Round.GetNextRound()
	}

	if gs.isHandOver() {
		gs.IsHandActive = false
		return true
	}
	return false
}

func (gs *State) processPots() {
	var processPotsPQ ProcessPotsPQ
	totalRoundPot := 0
	highestBetAmountPlayerIndex := -1
	highestBetAmount := 0
	secondHighestBetAmount := 0
	isHighestBetCalled := false
	for i, player := range gs.Players {
		if player == nil {
			continue
		}
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
		amountToRefund := highestBetAmount - secondHighestBetAmount
		gs.Players[highestBetAmountPlayerIndex].Stack += amountToRefund
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
			if player == nil {
				continue
			}
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

func (gs *State) processEndGame() {

	if len(gs.Pots) != len(gs.PotContenders) {
		log.Fatal("bug: Pots and PotContenders are out of sync")
	}

	// Go through the side pots from latest to earliest
	for potIndex := len(gs.Pots) - 1; potIndex >= 0; potIndex-- {
		if onePlayerStanding, player := gs.isOnePlayerStanding(gs.PotContenders[potIndex]); onePlayerStanding {
			gs.Players[player].Stack += gs.Pots[potIndex]
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

			numberOfWinners := 1
			for i := 1; i < len(rankedHands) && CompareStrengths(rankedHands[0].HandStrength, rankedHands[i].HandStrength) == 0; i++ {
				numberOfWinners++
			}

			for i := 0; i < numberOfWinners; i++ {
				gs.Players[rankedHands[i].PlayerIndex].Stack += gs.Pots[potIndex] / numberOfWinners
			}
			//consequence.PotRemainder += gs.Pots[potIndex] % numberOfWinners
		}
	}
}

// Returns index of player that is `n` active players to the right of `base`
func (gs *State) getNActivePlayerIndexFromIndex(base, n int) int {
	index := base
	count := 0
	for count != n {
		index = (index + 1) % 9
		for gs.Players[index] == nil || gs.Players[index].Folded {
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
	playerStanding := -1
	for i, player := range gs.Players {
		if player != nil && !player.Folded && containsIntInIntSlice(playersToConsider, i) {
			playersInHand++
			playerStanding = i
		}
	}
	return playersInHand == 1, playerStanding
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
