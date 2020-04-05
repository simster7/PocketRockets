package engine

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var shuffler = NewDeterministicShuffler()

func getNewTestGameState(seats map[int]int) State {
	players := new(Players)
	for seat, stack := range seats {
		players[seat] = &Player{Name: string(seat), Stack: stack}
	}
	return NewHandState(*players, 0)
}

func TestProcessPots(t *testing.T) {

	// Test standard pot
	gs := getNewTestGameState(map[int]int{2: 100, 5: 100, 7: 100})
	gs.Players[2].Bet = 10
	gs.Players[5].Bet = 10
	gs.Players[7].Bet = 10

	gs.Pots = []int{0}
	gs.PotContenders = [][]int{AllPlayers}

	gs.processPots()

	assert.Equal(t, []int{30}, gs.Pots)
	assert.Equal(t, [][]int{{0, 1, 2, 3, 4, 5, 6, 7, 8}}, gs.PotContenders)

	// Test one all in
	gs = getNewTestGameState(map[int]int{2: 100, 5: 100, 7: 100})
	gs.Players[2].Bet = 30
	gs.Players[5].Bet = 30
	gs.Players[7].Bet = 20
	gs.Players[7].IsAllIn = true

	gs.Pots = []int{0}
	gs.PotContenders = [][]int{AllPlayers}

	gs.processPots()

	assert.Equal(t, []int{60, 20}, gs.Pots)
	assert.Equal(t, [][]int{{0, 1, 2, 3, 4, 5, 6, 7, 8}, {0, 1, 2, 3, 4, 5, 6, 8}}, gs.PotContenders)

	// Test two all ins
	gs = getNewTestGameState(map[int]int{0: 100, 2: 100, 5: 100, 7: 100})
	gs.Players[0].Bet = 40
	gs.Players[2].Bet = 40
	gs.Players[5].Bet = 30
	gs.Players[5].IsAllIn = true
	gs.Players[7].Bet = 20
	gs.Players[7].IsAllIn = true

	gs.Pots = []int{0}
	gs.PotContenders = [][]int{AllPlayers}

	gs.processPots()

	assert.Equal(t, []int{80, 30, 20}, gs.Pots)
	assert.Equal(t, [][]int{{0, 1, 2, 3, 4, 5, 6, 7, 8}, {0, 1, 2, 3, 4, 5, 6, 8}, {0, 1, 2, 3, 4, 6, 8}}, gs.PotContenders)

	// Test two all ins with folds
	gs = getNewTestGameState(map[int]int{0: 100, 2: 100, 3: 100, 5: 100, 7: 100})
	gs.Players[0].Bet = 40
	gs.Players[2].Bet = 40
	gs.Players[3].Bet = 10
	gs.Players[5].Bet = 30
	gs.Players[5].IsAllIn = true
	gs.Players[7].Bet = 20
	gs.Players[7].IsAllIn = true

	gs.Pots = []int{0}
	gs.PotContenders = [][]int{AllPlayers}

	gs.processPots()

	assert.Equal(t, []int{90, 30, 20}, gs.Pots)
	assert.Equal(t, [][]int{{0, 1, 2, 3, 4, 5, 6, 7, 8}, {0, 1, 2, 3, 4, 5, 6, 8}, {0, 1, 2, 3, 4, 6, 8}}, gs.PotContenders)

	// Test all in with over-bet amount that can't be matched
	gs = getNewTestGameState(map[int]int{5: 100, 7: 100})
	gs.Players[5].Bet = 30
	gs.Players[7].Bet = 10
	gs.Players[7].IsAllIn = true

	gs.Pots = []int{0}
	gs.PotContenders = [][]int{AllPlayers}

	gs.processPots()

	assert.Equal(t, []int{20, 0}, gs.Pots)
	assert.Equal(t, [][]int{AllPlayers, {0, 1, 2, 3, 4, 5, 6, 8}}, gs.PotContenders)
}

func TestGameBasicSplitPot(t *testing.T) {
	gs := getNewTestGameState(map[int]int{0: 100, 1: 100, 3: 100, 4: 100, 5: 100, 6: 100})
	gs.DealHand(2, 1, shuffler.Shuffle())

	// Pre flop
	assert.Equal(t, RoundPreFlop, gs.Round)
	assert.Len(t, gs.getCommunityCards(), 0)

	err := gs.TakeAction(Action{ActionType: ActionTypeBet, Value: 5})
	assert.NoError(t, err)
	assert.Equal(t, 93, gs.Players[5].Stack)
	err = gs.TakeAction(Action{ActionType: ActionTypeCall})
	assert.NoError(t, err)
	assert.Equal(t, 93, gs.Players[6].Stack)
	err = gs.TakeAction(Action{ActionType: ActionTypeFold})
	assert.NoError(t, err)
	assert.Equal(t, 100, gs.Players[0].Stack)
	assert.True(t, gs.Players[0].Folded)
	err = gs.TakeAction(Action{ActionType: ActionTypeCall})
	assert.NoError(t, err)
	assert.Equal(t, 93, gs.Players[1].Stack)
	err = gs.TakeAction(Action{ActionType: ActionTypeCall})
	assert.NoError(t, err)
	assert.Equal(t, 93, gs.Players[3].Stack)
	err = gs.TakeAction(Action{ActionType: ActionTypeFold})
	assert.NoError(t, err)
	assert.True(t, gs.Players[4].Folded)

	// Flop
	assert.Equal(t, RoundFlop, gs.Round)
	assert.Len(t, gs.getCommunityCards(), 3)
	assert.Equal(t, 30, gs.Pots[0])

	// Can't call at start of round
	err = gs.TakeAction(Action{ActionType: ActionTypeCall})
	assert.Error(t, err)

	err = gs.TakeAction(Action{ActionType: ActionTypeCheck})
	assert.NoError(t, err)
	err = gs.TakeAction(Action{ActionType: ActionTypeCheck})
	assert.NoError(t, err)
	err = gs.TakeAction(Action{ActionType: ActionTypeBet, Value: 10})
	assert.NoError(t, err)
	assert.Equal(t, 83, gs.Players[6].Stack)
	// Can't check a bet
	err = gs.TakeAction(Action{ActionType: ActionTypeCheck})
	assert.Error(t, err)
	err = gs.TakeAction(Action{ActionType: ActionTypeFold})
	assert.NoError(t, err)
	assert.True(t, gs.Players[1].Folded)
	err = gs.TakeAction(Action{ActionType: ActionTypeCall})
	assert.NoError(t, err)
	assert.Equal(t, 83, gs.Players[3].Stack)
	err = gs.TakeAction(Action{ActionType: ActionTypeCall})
	assert.NoError(t, err)
	assert.Equal(t, 83, gs.Players[5].Stack)

	// Turn
	assert.Equal(t, RoundTurn, gs.Round)
	assert.Len(t, gs.getCommunityCards(), 4)
	assert.Equal(t, 60, gs.Pots[0])

	err = gs.TakeAction(Action{ActionType: ActionTypeCheck})
	assert.NoError(t, err)
	err = gs.TakeAction(Action{ActionType: ActionTypeCheck})
	assert.NoError(t, err)
	err = gs.TakeAction(Action{ActionType: ActionTypeCheck})
	assert.NoError(t, err)

	// River
	assert.Equal(t, RoundRiver, gs.Round)
	assert.Len(t, gs.getCommunityCards(), 5)
	assert.Equal(t, 60, gs.Pots[0])

	err = gs.TakeAction(Action{ActionType: ActionTypeBet, Value: 10})
	assert.NoError(t, err)
	assert.Equal(t, 73, gs.Players[3].Stack)
	err = gs.TakeAction(Action{ActionType: ActionTypeFold})
	assert.NoError(t, err)
	assert.True(t, gs.Players[5].Folded)
	err = gs.TakeAction(Action{ActionType: ActionTypeCall})
	assert.NoError(t, err)

	// Post River
	// Board hits a flush, split pot
	assert.Equal(t, RoundHandEnd, gs.Round)
	assert.Len(t, gs.getCommunityCards(), 5)
	assert.Equal(t, 80, gs.Pots[0])
	assert.False(t, gs.IsHandActive)
	assert.Equal(t, 113, gs.Players[3].Stack)
	assert.Equal(t, 113, gs.Players[6].Stack)
}

func TestGameMultiround(t *testing.T) {
	gs := getNewTestGameState(map[int]int{2: 100, 5: 100, 7: 100})
	gs.DealHand(2, 1, shuffler.Shuffle())

	// Pre flop
	assert.Equal(t, RoundPreFlop, gs.Round)
	err := gs.TakeAction(Action{ActionType: ActionTypeBet, Value: 8})
	assert.NoError(t, err)
	assert.Equal(t, 90, gs.Players[2].Stack)
	// Can't check a call
	err = gs.TakeAction(Action{ActionType: ActionTypeCheck})
	assert.Error(t, err)
	err = gs.TakeAction(Action{ActionType: ActionTypeCall})
	assert.NoError(t, err)
	assert.Equal(t, 90, gs.Players[2].Stack)
	err = gs.TakeAction(Action{ActionType: ActionTypeCall})
	assert.NoError(t, err)

	// Flop
	assert.Equal(t, RoundFlop, gs.Round)
	err = gs.TakeAction(Action{ActionType: ActionTypeCheck})
	assert.NoError(t, err)
	err = gs.TakeAction(Action{ActionType: ActionTypeCheck})
	assert.NoError(t, err)
	err = gs.TakeAction(Action{ActionType: ActionTypeCheck})
	assert.NoError(t, err)

	// Turn
	assert.Equal(t, RoundTurn, gs.Round)
	// Can't bet more than stack
	err = gs.TakeAction(Action{ActionType: ActionTypeBet, Value: 1000})
	assert.Error(t, err)
	err = gs.TakeAction(Action{ActionType: ActionTypeBet, Value: 10})
	assert.NoError(t, err)
	assert.Equal(t, 80, gs.Players[5].Stack)
	err = gs.TakeAction(Action{ActionType: ActionTypeCall})
	assert.NoError(t, err)
	assert.Equal(t, 80, gs.Players[7].Stack)
	err = gs.TakeAction(Action{ActionType: ActionTypeFold})
	assert.NoError(t, err)

	// River
	assert.Equal(t, RoundRiver, gs.Round)
	err = gs.TakeAction(Action{ActionType: ActionTypeCheck})
	assert.NoError(t, err)
	err = gs.TakeAction(Action{ActionType: ActionTypeCheck})
	assert.NoError(t, err)

	assert.Equal(t, RoundHandEnd, gs.Round)
	assert.Equal(t, 90, gs.Players[2].Stack)
	assert.Equal(t, 80, gs.Players[5].Stack)
	assert.Equal(t, 130, gs.Players[7].Stack)

	assert.False(t, gs.IsHandActive)

	//gs = getNewTestGameState(map[int]int{2: 90, 5: 80, 7: 130, 8: 100})
	gs.Players[8] = &Player{Name: "8", Stack: 100}
	gs.DealHand(2, 1, shuffler.Shuffle())

	// Pre flop
	assert.Equal(t, RoundPreFlop, gs.Round)
	assert.Len(t, gs.getCommunityCards(), 0)

	err = gs.TakeAction(Action{ActionType: ActionTypeCall})
	assert.NoError(t, err)
	err = gs.TakeAction(Action{ActionType: ActionTypeBet, Value: 10})
	assert.NoError(t, err)
	err = gs.TakeAction(Action{ActionType: ActionTypeFold})
	assert.NoError(t, err)
	err = gs.TakeAction(Action{ActionType: ActionTypeCall})
	assert.NoError(t, err)
	err = gs.TakeAction(Action{ActionType: ActionTypeCall})
	assert.NoError(t, err)

	// Flop 37
	assert.Equal(t, RoundFlop, gs.Round)
	assert.Len(t, gs.getCommunityCards(), 3)

	err = gs.TakeAction(Action{ActionType: ActionTypeCheck})
	assert.NoError(t, err)
	err = gs.TakeAction(Action{ActionType: ActionTypeCheck})
	assert.NoError(t, err)
	err = gs.TakeAction(Action{ActionType: ActionTypeBet, Value: 10})
	assert.NoError(t, err)
	err = gs.TakeAction(Action{ActionType: ActionTypeFold})
	assert.NoError(t, err)
	assert.True(t, gs.Players[8].Folded)
	err = gs.TakeAction(Action{ActionType: ActionTypeCall})
	assert.NoError(t, err)

	// Turn
	assert.Equal(t, RoundTurn, gs.Round)
	assert.Len(t, gs.getCommunityCards(), 4)

	err = gs.TakeAction(Action{ActionType: ActionTypeCheck})
	assert.NoError(t, err)
	err = gs.TakeAction(Action{ActionType: ActionTypeCheck})
	assert.NoError(t, err)

	// River
	assert.Equal(t, RoundRiver, gs.Round)
	assert.Len(t, gs.getCommunityCards(), 5)

	err = gs.TakeAction(Action{ActionType: ActionTypeCheck})
	assert.NoError(t, err)
	err = gs.TakeAction(Action{ActionType: ActionTypeCheck})
	assert.NoError(t, err)

	// PostRiver
	assert.Equal(t, RoundHandEnd, gs.Round)

	assert.Equal(t, 115, gs.Players[5].Stack)
	assert.Equal(t, 68, gs.Players[2].Stack)
	assert.Equal(t, 88, gs.Players[8].Stack)
	assert.Equal(t, 129, gs.Players[7].Stack)

	//gs = getNewTestGameState(map[int]int{2: 68, 5: 115, 7: 129, 8: 88})
	gs.DealHand(2, 1, shuffler.Shuffle())

	// Pre flop
	assert.Equal(t, RoundPreFlop, gs.Round)
	assert.Len(t, gs.getCommunityCards(), 0)

	err = gs.TakeAction(Action{ActionType: ActionTypeBet, Value: 10})
	assert.NoError(t, err)
	err = gs.TakeAction(Action{ActionType: ActionTypeCall})
	assert.NoError(t, err)
	err = gs.TakeAction(Action{ActionType: ActionTypeCall})
	assert.NoError(t, err)
	err = gs.TakeAction(Action{ActionType: ActionTypeCall})
	assert.NoError(t, err)

	assert.Equal(t, 103, gs.Players[5].Stack)
	assert.Equal(t, 56, gs.Players[2].Stack)
	assert.Equal(t, 76, gs.Players[8].Stack)
	assert.Equal(t, 117, gs.Players[7].Stack)

	// Pre flop
	assert.Equal(t, RoundFlop, gs.Round)
	// DONK!
	err = gs.TakeAction(Action{ActionType: ActionTypeBet, Value: 50})
	assert.NoError(t, err)
	assert.Equal(t, 26, gs.Players[8].Stack)
	err = gs.TakeAction(Action{ActionType: ActionTypeFold})
	assert.NoError(t, err)
	err = gs.TakeAction(Action{ActionType: ActionTypeFold})
	assert.NoError(t, err)
	err = gs.TakeAction(Action{ActionType: ActionTypeFold})
	assert.NoError(t, err)

	assert.False(t, gs.IsHandActive)
	assert.Equal(t, 103, gs.Players[5].Stack)
	assert.Equal(t, 56, gs.Players[2].Stack)
	assert.Equal(t, 124, gs.Players[8].Stack)
	assert.Equal(t, 117, gs.Players[7].Stack)
}

func TestGameAllInSimple(t *testing.T) {
	gs := getNewTestGameState(map[int]int{2: 20, 5: 50, 7: 100})

	// Pre flop
	// Bet is 10 each, 2 is left with 10 at round end
	assert.Equal(t, RoundPreFlop, gs.Round)
	err := gs.TakeAction(Action{ActionType: ActionTypeBet, Value: 8})
	assert.NoError(t, err)
	assert.Equal(t, 10, gs.Players[2].Stack)
	err = gs.TakeAction(Action{ActionType: ActionTypeCall})
	assert.NoError(t, err)
	assert.Equal(t, 40, gs.Players[5].Stack)
	err = gs.TakeAction(Action{ActionType: ActionTypeCall})
	assert.NoError(t, err)
	assert.Equal(t, 90, gs.Players[7].Stack)

	// Flop
	// Bet is 30 each, 2 can only afford 10. Main pot becomes 60 (30 from flop + 10 (2's all-in) * 3 (players))
	// Side pot becomes 40 (20 each from 5 and 7)
	assert.Equal(t, RoundFlop, gs.Round)
	err = gs.TakeAction(Action{ActionType: ActionTypeBet, Value: 30})
	assert.NoError(t, err)
	assert.Equal(t, 10, gs.Players[5].Stack)
	err = gs.TakeAction(Action{ActionType: ActionTypeCall})
	assert.NoError(t, err)
	assert.Equal(t, 60, gs.Players[7].Stack)
	err = gs.TakeAction(Action{ActionType: ActionTypeCall})
	assert.NoError(t, err)
	assert.Equal(t, 0, gs.Players[2].Stack)

	// Turn
	assert.Equal(t, RoundTurn, gs.Round)
	err = gs.TakeAction(Action{ActionType: ActionTypeCheck})
	assert.NoError(t, err)
	err = gs.TakeAction(Action{ActionType: ActionTypeCheck})
	assert.NoError(t, err)

	// River
	assert.Equal(t, RoundRiver, gs.Round)
	err = gs.TakeAction(Action{ActionType: ActionTypeCheck})
	assert.NoError(t, err)
	err = gs.TakeAction(Action{ActionType: ActionTypeCheck})
	assert.NoError(t, err)

	// HandEnd
	// 2 wins main pot of 60, 7 wins side pot of 40
	assert.Equal(t, RoundHandEnd, gs.Round)
	assert.Equal(t, 60, gs.Players[2].Stack)
	assert.Equal(t, 100, gs.Players[7].Stack)
	assert.Equal(t, 10, gs.Players[5].Stack)
}

func TestGameAllInTwoSidePots(t *testing.T) {
	gs := getNewTestGameState(map[int]int{2: 20, 5: 50, 7: 100, 8: 30})

	// Pre flop
	// Bet is 10 each, 2 is left with 10 at round end and Jarry with 20
	assert.Equal(t, RoundPreFlop, gs.Round)
	err := gs.TakeAction(Action{ActionType: ActionTypeBet, Value: 8})
	assert.NoError(t, err)
	assert.Equal(t, 20, gs.Players[2].Stack)
	err = gs.TakeAction(Action{ActionType: ActionTypeCall})
	assert.NoError(t, err)
	assert.Equal(t, 10, gs.Players[2].Stack)
	err = gs.TakeAction(Action{ActionType: ActionTypeCall})
	assert.NoError(t, err)
	assert.Equal(t, 40, gs.Players[5].Stack)
	err = gs.TakeAction(Action{ActionType: ActionTypeCall})
	assert.NoError(t, err)
	assert.Equal(t, 90, gs.Players[7].Stack)

	// Flop
	// Bet is 15 each, 2 can only afford 10. Main pot becomes 80 (40 from preflop + 10 (2's all-in) * 4 (players))
	// First side pot becomes 15 (5 each from 5, 7, and 8)
	assert.Equal(t, RoundFlop, gs.Round)
	err = gs.TakeAction(Action{ActionType: ActionTypeBet, Value: 15})
	assert.NoError(t, err)
	assert.Equal(t, 25, gs.Players[5].Stack)
	err = gs.TakeAction(Action{ActionType: ActionTypeCall})
	assert.NoError(t, err)
	assert.Equal(t, 75, gs.Players[7].Stack)
	err = gs.TakeAction(Action{ActionType: ActionTypeCall})
	assert.NoError(t, err)
	assert.Equal(t, 5, gs.Players[8].Stack)
	err = gs.TakeAction(Action{ActionType: ActionTypeCall})
	assert.NoError(t, err)
	assert.Equal(t, 0, gs.Players[2].Stack)

	// Turn
	// Bet is 15 each, 8 can only afford 5. First side pot becomes 30 (15 from flop + 5 (8's all-in) * 3 (players))
	// Second side pot becomes 20 (10 each from 5 and 7)
	assert.Equal(t, RoundTurn, gs.Round)
	err = gs.TakeAction(Action{ActionType: ActionTypeBet, Value: 15})
	assert.NoError(t, err)
	assert.Equal(t, 10, gs.Players[5].Stack)
	err = gs.TakeAction(Action{ActionType: ActionTypeCall})
	assert.NoError(t, err)
	assert.Equal(t, 60, gs.Players[7].Stack)
	err = gs.TakeAction(Action{ActionType: ActionTypeCall})
	assert.NoError(t, err)
	assert.Equal(t, 0, gs.Players[8].Stack)
	// Jason is already all-in
	err = gs.TakeAction(Action{ActionType: ActionTypeCall})
	assert.Error(t, err)

	// River
	assert.Equal(t, RoundRiver, gs.Round)
	err = gs.TakeAction(Action{ActionType: ActionTypeCheck})
	assert.NoError(t, err)
	err = gs.TakeAction(Action{ActionType: ActionTypeCheck})
	assert.NoError(t, err)

	// HandEnd
	// First showdown is 5 vs 7. 7 wins second side pot of 20 and now has 80
	// Second showdown is 5 vs 7 vs 8. 8 wins first sidepot of 30 and now has 30
	// Last showdown is family pot. 2 wins main pot of 80 and now has 80
	assert.Equal(t, RoundHandEnd, gs.Round)
	assert.Equal(t, 80, gs.Players[2].Stack)
	assert.Equal(t, 80, gs.Players[7].Stack)
	assert.Equal(t, 10, gs.Players[5].Stack)
	assert.Equal(t, 30, gs.Players[8].Stack)

}

func TestGameAllInWithFold(t *testing.T) {
	gs := getNewTestGameState(map[int]int{2: 20, 5: 50, 7: 100, 8: 30})

	// Pre flop
	// Bet is 10 each, 2 is left with 10 at round end and 8 with 20
	assert.Equal(t, RoundPreFlop, gs.Round)
	err := gs.TakeAction(Action{ActionType: ActionTypeBet, Value: 8})
	assert.NoError(t, err)
	assert.Equal(t, 20, gs.Players[2].Stack)
	err = gs.TakeAction(Action{ActionType: ActionTypeCall})
	assert.NoError(t, err)
	assert.Equal(t, 10, gs.Players[2].Stack)
	err = gs.TakeAction(Action{ActionType: ActionTypeCall})
	assert.NoError(t, err)
	assert.Equal(t, 40, gs.Players[5].Stack)
	err = gs.TakeAction(Action{ActionType: ActionTypeCall})
	assert.NoError(t, err)
	assert.Equal(t, 90, gs.Players[7].Stack)

	// Flop
	// Main pot becomes 75 (40 from preflop + 10 (2's all-in) * 3 players + 5 (8's bet-fold))
	// First side pot becomes 10 (5 each from 5, 8)
	assert.Equal(t, RoundFlop, gs.Round)
	err = gs.TakeAction(Action{ActionType: ActionTypeCheck})
	assert.NoError(t, err)
	assert.Equal(t, 40, gs.Players[5].Stack)
	err = gs.TakeAction(Action{ActionType: ActionTypeBet, Value: 5})
	assert.NoError(t, err)
	assert.Equal(t, 85, gs.Players[7].Stack)
	// Calls 5, bets 10. Total: 15
	err = gs.TakeAction(Action{ActionType: ActionTypeBet, Value: 10})
	assert.NoError(t, err)
	assert.Equal(t, 5, gs.Players[8].Stack)
	err = gs.TakeAction(Action{ActionType: ActionTypeCall})
	assert.NoError(t, err)
	assert.Equal(t, 0, gs.Players[2].Stack)
	err = gs.TakeAction(Action{ActionType: ActionTypeCall})
	assert.NoError(t, err)
	assert.Equal(t, 25, gs.Players[5].Stack)
	err = gs.TakeAction(Action{ActionType: ActionTypeFold})
	assert.NoError(t, err)

	// Turn
	// Bet is 15 each, 8 can only afford 5, so side pot becomes 20 (10 from flop + 5 (8's all-in) * 2 (players))
	assert.Equal(t, RoundTurn, gs.Round)
	err = gs.TakeAction(Action{ActionType: ActionTypeBet, Value: 15})
	assert.NoError(t, err)
	assert.Equal(t, 10, gs.Players[5].Stack)

	err = gs.TakeAction(Action{ActionType: ActionTypeCall})
	assert.NoError(t, err)
	assert.Equal(t, 0, gs.Players[8].Stack)
	// 2 is already all-in
	err = gs.TakeAction(Action{ActionType: ActionTypeCall})
	assert.Error(t, err)

	// 5 bet 15, but only 5 was called, so 10 should have been added back to 5's stack
	assert.Equal(t, 20, gs.Players[5].Stack)

	// River
	// TODO: No more action, this should be auto done
	assert.Equal(t, RoundRiver, gs.Round)
	err = gs.TakeAction(Action{ActionType: ActionTypeCheck})
	assert.NoError(t, err)

	// HandEnd
	// First showdown is 5 vs 8. 8 wins side pot of 20 and now has 20
	// Second showdown is 5 vs 8 vs 2. 2 wins main pot of 75 and now has 75
	assert.Equal(t, RoundHandEnd, gs.Round)
	assert.Equal(t, 75, gs.Players[2].Stack)
	assert.Equal(t, 85, gs.Players[7].Stack)
	assert.Equal(t, 20, gs.Players[5].Stack)
	assert.Equal(t, 20, gs.Players[8].Stack)
}

func TestGamePreFlopOption(t *testing.T) {
	gs := getNewTestGameState(map[int]int{2: 100, 5: 100, 7: 100})

	assert.Equal(t, RoundPreFlop, gs.Round)
	err := gs.TakeAction(Action{ActionType: ActionTypeCall})
	assert.NoError(t, err)
	err = gs.TakeAction(Action{ActionType: ActionTypeCall})
	assert.NoError(t, err)
	// Option check
	err = gs.TakeAction(Action{ActionType: ActionTypeCheck})
	assert.NoError(t, err)

	assert.Equal(t, RoundFlop, gs.Round)
	err = gs.TakeAction(Action{ActionType: ActionTypeCheck})
	assert.NoError(t, err)
	assert.Equal(t, 98, gs.Players[2].Stack)
	assert.Equal(t, 98, gs.Players[5].Stack)
	assert.Equal(t, 98, gs.Players[7].Stack)

	gs = getNewTestGameState(map[int]int{2: 100, 5: 100, 7: 100})

	assert.Equal(t, RoundPreFlop, gs.Round)
	err = gs.TakeAction(Action{ActionType: ActionTypeCall})
	assert.NoError(t, err)
	err = gs.TakeAction(Action{ActionType: ActionTypeCall})
	assert.NoError(t, err)
	// Option raise
	err = gs.TakeAction(Action{ActionType: ActionTypeBet, Value: 5})
	assert.NoError(t, err)
	// Still pre flop
	assert.Equal(t, RoundPreFlop, gs.Round)
	err = gs.TakeAction(Action{ActionType: ActionTypeCall})
	assert.NoError(t, err)
	err = gs.TakeAction(Action{ActionType: ActionTypeCall})
	assert.NoError(t, err)

	assert.Equal(t, RoundFlop, gs.Round)
	err = gs.TakeAction(Action{ActionType: ActionTypeCheck})
	assert.NoError(t, err)
	assert.Equal(t, 93, gs.Players[2].Stack)
	assert.Equal(t, 93, gs.Players[5].Stack)
	assert.Equal(t, 93, gs.Players[7].Stack)
}

