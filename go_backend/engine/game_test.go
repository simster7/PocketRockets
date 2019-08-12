package engine

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGameBasicSplitPot(t *testing.T) {
	// TODO Replace GameState calls with API calls
	game := NewDeterministicGame(1, 2, getDeck)
	grace := NewPlayer("Grace", 100)
	err := game.SitPlayer(&grace, 0)
	assert.NoError(t, err)
	jason := NewPlayer("Jason", 100)
	err = game.SitPlayer(&jason, 1)
	assert.NoError(t, err)
	simon := NewPlayer("Simon", 100)
	err = game.SitPlayer(&simon, 3)
	assert.NoError(t, err)
	hersh := NewPlayer("Hersh", 100)
	err = game.SitPlayer(&hersh, 4)
	assert.NoError(t, err)
	chien := NewPlayer("Chien", 100)
	err = game.SitPlayer(&chien, 5)
	assert.NoError(t, err)
	jarry := NewPlayer("Jarry", 100)
	err = game.SitPlayer(&jarry, 6)
	assert.NoError(t, err)

	game.DealHand()

	// Pre flop
	assert.Equal(t, PreFlop, game.GameState.Round)
	assert.Nil(t, game.GameState.getCommunityCards())
	// Can't play out of turn
	err = game.TakeAction(&jarry, Action{ActionType: Bet, Value: 5})
	assert.Error(t, err)
	assert.Equal(t, 100, jarry.Stack)
	err = game.TakeAction(&chien, Action{ActionType: Bet, Value: 5})
	assert.NoError(t, err)
	assert.Equal(t, 93, chien.Stack)
	err = game.TakeAction(&jarry, Action{ActionType: Call})
	assert.NoError(t, err)
	assert.Equal(t, 93, jarry.Stack)
	err = game.TakeAction(&grace, Action{ActionType: Fold})
	assert.NoError(t, err)
	assert.Equal(t, 100, grace.Stack)
	assert.Equal(t, true, grace.Folded)
	err = game.TakeAction(&jason, Action{ActionType: Call})
	assert.NoError(t, err)
	err = game.TakeAction(&simon, Action{ActionType: Call})
	assert.NoError(t, err)
	err = game.TakeAction(&hersh, Action{ActionType: Fold})
	assert.NoError(t, err)
	assert.Equal(t, true, hersh.Folded)

	// Flop
	assert.Equal(t, Flop, game.GameState.Round)
	assert.NotNil(t, game.GameState.getCommunityCards())
	assert.Len(t, game.GameState.getCommunityCards(), 3)
	assert.Equal(t, game.GameState.Pots[0], 30)

	// Can't call at start of round
	err = game.TakeAction(&simon, Action{ActionType: Call})
	assert.Error(t, err)
	// Can't play when folded
	err = game.TakeAction(&hersh, Action{ActionType: Call})
	assert.Error(t, err)
	err = game.TakeAction(&simon, Action{ActionType: Check})
	assert.NoError(t, err)
	err = game.TakeAction(&chien, Action{ActionType: Check})
	assert.NoError(t, err)
	err = game.TakeAction(&jarry, Action{ActionType: Bet, Value: 10})
	assert.NoError(t, err)
	assert.Equal(t, 83, jarry.Stack)
	// Can't check a bet
	err = game.TakeAction(&jason, Action{ActionType: Check})
	assert.Error(t, err)
	err = game.TakeAction(&jason, Action{ActionType: Fold})
	assert.NoError(t, err)
	assert.Equal(t, true, jason.Folded)
	err = game.TakeAction(&simon, Action{ActionType: Call})
	assert.NoError(t, err)
	assert.Equal(t, 83, simon.Stack)
	err = game.TakeAction(&chien, Action{ActionType: Call})
	assert.NoError(t, err)
	assert.Equal(t, 83, chien.Stack)

	// Turn
	assert.Equal(t, Turn, game.GameState.Round)
	assert.Len(t, game.GameState.getCommunityCards(), 4)
	assert.Equal(t, game.GameState.Pots[0], 60)

	err = game.TakeAction(&simon, Action{ActionType: Check})
	assert.NoError(t, err)
	err = game.TakeAction(&chien, Action{ActionType: Check})
	assert.NoError(t, err)
	err = game.TakeAction(&jarry, Action{ActionType: Check})
	assert.NoError(t, err)

	// River
	assert.Equal(t, River, game.GameState.Round)
	assert.Len(t, game.GameState.getCommunityCards(), 5)
	assert.Equal(t, game.GameState.Pots[0], 60)

	err = game.TakeAction(&simon, Action{ActionType: Bet, Value: 10})
	assert.NoError(t, err)
	assert.Equal(t, 73, simon.Stack)
	err = game.TakeAction(&chien, Action{ActionType: Fold})
	assert.NoError(t, err)
	assert.Equal(t, true, chien.Folded)
	err = game.TakeAction(&jarry, Action{ActionType: Call})
	assert.NoError(t, err)

	// Post River
	// Board hits a flush, split pot
	assert.Equal(t, HandEnd, game.GameState.Round)
	assert.Equal(t, game.GameState.Pots[0], 80)
	assert.Len(t, game.GameState.getCommunityCards(), 5)
	assert.False(t, game.IsHandActive)
	assert.False(t, game.GameState.IsHandActive)
	assert.Equal(t, 113, simon.Stack)
	assert.Equal(t, 113, jarry.Stack)
}

func TestGameMultiround(t *testing.T) {
	game := NewDeterministicGame(1, 2, getDeck)
	jason := NewPlayer("Jason", 100)
	err := game.SitPlayer(&jason, 2)
	assert.NoError(t, err)
	simon := NewPlayer("Simon", 100)
	err = game.SitPlayer(&simon, 5)
	assert.NoError(t, err)
	chien := NewPlayer("Chien", 100)
	err = game.SitPlayer(&chien, 7)
	assert.NoError(t, err)

	game.DealHand()

	// Pre flop
	assert.Equal(t, PreFlop, game.GameState.Round)
	err = game.TakeAction(&jason, Action{ActionType: Bet, Value: 8})
	assert.NoError(t, err)
	assert.Equal(t, 90, jason.Stack)
	// Can't check a call
	err = game.TakeAction(&simon, Action{ActionType: Check})
	assert.Error(t, err)
	err = game.TakeAction(&simon, Action{ActionType: Call})
	assert.NoError(t, err)
	assert.Equal(t, 90, jason.Stack)
	err = game.TakeAction(&chien, Action{ActionType: Call})
	assert.NoError(t, err)

	// Flop
	assert.Equal(t, Flop, game.GameState.Round)
	err = game.TakeAction(&simon, Action{ActionType: Check})
	assert.NoError(t, err)
	err = game.TakeAction(&chien, Action{ActionType: Check})
	assert.NoError(t, err)
	err = game.TakeAction(&jason, Action{ActionType: Check})
	assert.NoError(t, err)

	// Turn
	assert.Equal(t, Turn, game.GameState.Round)
	// Can't bet more than stack
	err = game.TakeAction(&simon, Action{ActionType: Bet, Value: 1000})
	assert.Error(t, err)
	err = game.TakeAction(&simon, Action{ActionType: Bet, Value: 10})
	assert.NoError(t, err)
	assert.Equal(t, 80, simon.Stack)
	err = game.TakeAction(&chien, Action{ActionType: Call})
	assert.NoError(t, err)
	assert.Equal(t, 80, chien.Stack)
	err = game.TakeAction(&jason, Action{ActionType: Fold})
	assert.NoError(t, err)

	// River
	assert.Equal(t, River, game.GameState.Round)
	err = game.TakeAction(&simon, Action{ActionType: Check})
	assert.NoError(t, err)
	err = game.TakeAction(&chien, Action{ActionType: Check})
	assert.NoError(t, err)

	assert.Equal(t, HandEnd, game.GameState.Round)
	assert.Equal(t, 80, simon.Stack)
	assert.Equal(t, 130, chien.Stack)

	assert.False(t, game.IsHandActive)

	jarry := NewPlayer("Jarry", 100)
	err = game.SitPlayer(&jarry, 8)
	assert.NoError(t, err)

	game.DealHand()

	// Pre flop
	assert.Equal(t, PreFlop, game.GameState.Round)
	assert.Nil(t, game.GameState.getCommunityCards())

	// Can't play out of turn
	err = game.TakeAction(&simon, Action{ActionType: Bet, Value: 5})
	assert.Error(t, err)
	err = game.TakeAction(&chien, Action{ActionType: Bet, Value: 5})
	assert.Error(t, err)
	err = game.TakeAction(&jarry, Action{ActionType: Bet, Value: 5})
	assert.Error(t, err)

	err = game.TakeAction(&jason, Action{ActionType: Call})
	assert.NoError(t, err)
	err = game.TakeAction(&simon, Action{ActionType: Bet, Value: 10})
	assert.NoError(t, err)
	err = game.TakeAction(&chien, Action{ActionType: Fold})
	assert.NoError(t, err)
	err = game.TakeAction(&jarry, Action{ActionType: Call})
	assert.NoError(t, err)
	err = game.TakeAction(&jason, Action{ActionType: Call})
	assert.NoError(t, err)

	// Flop 37
	assert.Equal(t, Flop, game.GameState.Round)
	assert.Len(t, game.GameState.getCommunityCards(), 3)

	// Can't play out of turn
	err = game.TakeAction(&simon, Action{ActionType: Bet, Value: 5})
	assert.Error(t, err)
	err = game.TakeAction(&jason, Action{ActionType: Bet, Value: 5})
	assert.Error(t, err)

	err = game.TakeAction(&jarry, Action{ActionType: Check})
	assert.NoError(t, err)
	err = game.TakeAction(&jason, Action{ActionType: Check})
	assert.NoError(t, err)
	err = game.TakeAction(&simon, Action{ActionType: Bet, Value: 10})
	assert.NoError(t, err)
	err = game.TakeAction(&jarry, Action{ActionType: Fold})
	assert.NoError(t, err)
	assert.True(t, jarry.Folded)
	err = game.TakeAction(&jason, Action{ActionType: Call})
	assert.NoError(t, err)

	// Turn
	assert.Equal(t, Turn, game.GameState.Round)
	assert.Len(t, game.GameState.getCommunityCards(), 4)

	err = game.TakeAction(&jason, Action{ActionType: Check})
	assert.NoError(t, err)
	err = game.TakeAction(&simon, Action{ActionType: Check})
	assert.NoError(t, err)

	// River
	assert.Equal(t, River, game.GameState.Round)
	assert.Len(t, game.GameState.getCommunityCards(), 5)

	err = game.TakeAction(&jason, Action{ActionType: Check})
	assert.NoError(t, err)
	err = game.TakeAction(&simon, Action{ActionType: Check})
	assert.NoError(t, err)

	// PostRiver
	assert.Equal(t, HandEnd, game.GameState.Round)

	assert.Equal(t, 115, simon.Stack)
	assert.Equal(t, 68, jason.Stack)
	assert.Equal(t, 88, jarry.Stack)
	assert.Equal(t, 129, chien.Stack)

	game.DealHand()

	// Pre flop
	assert.Equal(t, PreFlop, game.GameState.Round)
	assert.Nil(t, game.GameState.getCommunityCards())

	err = game.TakeAction(&simon, Action{ActionType: Bet, Value: 10})
	assert.NoError(t, err)
	err = game.TakeAction(&chien, Action{ActionType: Call})
	assert.NoError(t, err)
	err = game.TakeAction(&jarry, Action{ActionType: Call})
	assert.NoError(t, err)
	err = game.TakeAction(&jason, Action{ActionType: Call})
	assert.NoError(t, err)

	assert.Equal(t, 103, simon.Stack)
	assert.Equal(t, 56, jason.Stack)
	assert.Equal(t, 76, jarry.Stack)
	assert.Equal(t, 117, chien.Stack)

	// Pre flop
	assert.Equal(t, Flop, game.GameState.Round)
	// DONK!
	err = game.TakeAction(&jarry, Action{ActionType: Bet, Value: 50})
	assert.NoError(t, err)
	assert.Equal(t, 26, jarry.Stack)
	err = game.TakeAction(&jason, Action{ActionType: Fold})
	assert.NoError(t, err)
	err = game.TakeAction(&simon, Action{ActionType: Fold})
	assert.NoError(t, err)
	err = game.TakeAction(&chien, Action{ActionType: Fold})
	assert.NoError(t, err)

	assert.False(t, game.IsHandActive)
	assert.Equal(t, 103, simon.Stack)
	assert.Equal(t, 56, jason.Stack)
	assert.Equal(t, 124, jarry.Stack)
	assert.Equal(t, 117, chien.Stack)
}

func TestGameAllInSimple(t *testing.T) {
	game := NewDeterministicGame(1, 2, getDeck)
	jason := NewPlayer("Jason", 20)
	err := game.SitPlayer(&jason, 2)
	assert.NoError(t, err)
	simon := NewPlayer("Simon", 50)
	err = game.SitPlayer(&simon, 5)
	assert.NoError(t, err)
	chien := NewPlayer("Chien", 100)
	err = game.SitPlayer(&chien, 7)
	assert.NoError(t, err)

	game.DealHand()

	// Pre flop
	// Bet is 10 each, Jason is left with 10 at round end
	assert.Equal(t, PreFlop, game.GameState.Round)
	err = game.TakeAction(&jason, Action{ActionType: Bet, Value: 8})
	assert.NoError(t, err)
	assert.Equal(t, 10, jason.Stack)
	err = game.TakeAction(&simon, Action{ActionType: Call})
	assert.NoError(t, err)
	assert.Equal(t, 40, simon.Stack)
	err = game.TakeAction(&chien, Action{ActionType: Call})
	assert.NoError(t, err)
	assert.Equal(t, 90, chien.Stack)

	// Flop
	// Bet is 30 each, Jason can only afford 10. Main pot becomes 60 (30 from flop + 10 (Jason's all-in) * 3 (players))
	// Side pot becomes 40 (20 each from Simon and Chien)
	assert.Equal(t, Flop, game.GameState.Round)
	err = game.TakeAction(&simon, Action{ActionType: Bet, Value: 30})
	assert.NoError(t, err)
	assert.Equal(t, 10, simon.Stack)
	err = game.TakeAction(&chien, Action{ActionType: Call})
	assert.NoError(t, err)
	assert.Equal(t, 60, chien.Stack)
	err = game.TakeAction(&jason, Action{ActionType: Call})
	assert.NoError(t, err)
	assert.Equal(t, 0, jason.Stack)

	// Turn
	assert.Equal(t, Turn, game.GameState.Round)
	err = game.TakeAction(&simon, Action{ActionType: Check})
	assert.NoError(t, err)
	err = game.TakeAction(&chien, Action{ActionType: Check})
	assert.NoError(t, err)

	// River
	assert.Equal(t, River, game.GameState.Round)
	err = game.TakeAction(&simon, Action{ActionType: Check})
	assert.NoError(t, err)
	err = game.TakeAction(&chien, Action{ActionType: Check})
	assert.NoError(t, err)

	// HandEnd
	// Jason wins main pot of 60, Chien wins side pot of 40
	assert.Equal(t, HandEnd, game.GameState.Round)
	assert.Equal(t, 60, jason.Stack)
	assert.Equal(t, 100, chien.Stack)
	assert.Equal(t, 10, simon.Stack)
}

func TestGameAllInTwoSidePots(t *testing.T) {
	game := NewDeterministicGame(1, 2, getDeck)
	jason := NewPlayer("Jason", 20)
	err := game.SitPlayer(&jason, 2)
	assert.NoError(t, err)
	simon := NewPlayer("Simon", 50)
	err = game.SitPlayer(&simon, 5)
	assert.NoError(t, err)
	chien := NewPlayer("Chien", 100)
	err = game.SitPlayer(&chien, 7)
	assert.NoError(t, err)
	jarry := NewPlayer("Jarry", 30)
	err = game.SitPlayer(&jarry, 8)
	assert.NoError(t, err)

	game.DealHand()

	// Pre flop
	// Bet is 10 each, Jason is left with 10 at round end and Jarry with 20
	assert.Equal(t, PreFlop, game.GameState.Round)
	err = game.TakeAction(&jarry, Action{ActionType: Bet, Value: 8})
	assert.NoError(t, err)
	assert.Equal(t, 20, jason.Stack)
	err = game.TakeAction(&jason, Action{ActionType: Call})
	assert.NoError(t, err)
	assert.Equal(t, 10, jason.Stack)
	err = game.TakeAction(&simon, Action{ActionType: Call})
	assert.NoError(t, err)
	assert.Equal(t, 40, simon.Stack)
	err = game.TakeAction(&chien, Action{ActionType: Call})
	assert.NoError(t, err)
	assert.Equal(t, 90, chien.Stack)

	// Flop
	// Bet is 15 each, Jason can only afford 10. Main pot becomes 80 (40 from preflop + 10 (Jason's all-in) * 4 (players))
	// First side pot becomes 15 (5 each from Simon, Chien, and Jarry)
	assert.Equal(t, Flop, game.GameState.Round)
	err = game.TakeAction(&simon, Action{ActionType: Bet, Value: 15})
	assert.NoError(t, err)
	assert.Equal(t, 25, simon.Stack)
	err = game.TakeAction(&chien, Action{ActionType: Call})
	assert.NoError(t, err)
	assert.Equal(t, 75, chien.Stack)
	err = game.TakeAction(&jarry, Action{ActionType: Call})
	assert.NoError(t, err)
	assert.Equal(t, 5, jarry.Stack)
	err = game.TakeAction(&jason, Action{ActionType: Call})
	assert.NoError(t, err)
	assert.Equal(t, 0, jason.Stack)

	// Turn
	// Bet is 15 each, Jarry can only afford 5. First side pot becomes 30 (15 from flop + 5 (Jarry's all-in) * 3 (players))
	// Second side pot becomes 20 (10 each from Simon and Chien)
	assert.Equal(t, Turn, game.GameState.Round)
	err = game.TakeAction(&simon, Action{ActionType: Bet, Value: 15})
	assert.NoError(t, err)
	assert.Equal(t, 10, simon.Stack)
	err = game.TakeAction(&chien, Action{ActionType: Call})
	assert.NoError(t, err)
	assert.Equal(t, 60, chien.Stack)
	err = game.TakeAction(&jarry, Action{ActionType: Call})
	assert.NoError(t, err)
	assert.Equal(t, 0, jarry.Stack)
	// Jason is already all-in
	err = game.TakeAction(&jason, Action{ActionType: Call})
	assert.Error(t, err)

	// River
	assert.Equal(t, River, game.GameState.Round)
	err = game.TakeAction(&simon, Action{ActionType: Check})
	assert.NoError(t, err)
	err = game.TakeAction(&chien, Action{ActionType: Check})
	assert.NoError(t, err)

	// HandEnd
	// First showdown is Simon vs Chien. Chien wins second side pot of 20 and now has 80
	// Second showdown is Simon vs Chien vs Jarry. Jarry wins first sidepot of 30 and now has 30
	// Last showdown is family pot. Jason wins main pot of 80 and now has 80
	assert.Equal(t, HandEnd, game.GameState.Round)
	assert.Equal(t, 80, jason.Stack)
	assert.Equal(t, 80, chien.Stack)
	assert.Equal(t, 10, simon.Stack)
	assert.Equal(t, 30, jarry.Stack)

}

func TestGameAllInWithFold(t *testing.T) {
	game := NewDeterministicGame(1, 2, getDeck)
	jason := NewPlayer("Jason", 20)
	err := game.SitPlayer(&jason, 2)
	assert.NoError(t, err)
	simon := NewPlayer("Simon", 50)
	err = game.SitPlayer(&simon, 5)
	assert.NoError(t, err)
	chien := NewPlayer("Chien", 100)
	err = game.SitPlayer(&chien, 7)
	assert.NoError(t, err)
	jarry := NewPlayer("Jarry", 30)
	err = game.SitPlayer(&jarry, 8)
	assert.NoError(t, err)

	game.DealHand()

	// Pre flop
	// Bet is 10 each, Jason is left with 10 at round end and Jarry with 20
	assert.Equal(t, PreFlop, game.GameState.Round)
	err = game.TakeAction(&jarry, Action{ActionType: Bet, Value: 8})
	assert.NoError(t, err)
	assert.Equal(t, 20, jason.Stack)
	err = game.TakeAction(&jason, Action{ActionType: Call})
	assert.NoError(t, err)
	assert.Equal(t, 10, jason.Stack)
	err = game.TakeAction(&simon, Action{ActionType: Call})
	assert.NoError(t, err)
	assert.Equal(t, 40, simon.Stack)
	err = game.TakeAction(&chien, Action{ActionType: Call})
	assert.NoError(t, err)
	assert.Equal(t, 90, chien.Stack)

	// Flop
	// Main pot becomes 75 (40 from preflop + 10 (Jason's all-in) * 3 players + 5 (chien's bet-fold))
	// First side pot becomes 10 (5 each from simon, jarry)
	assert.Equal(t, Flop, game.GameState.Round)
	err = game.TakeAction(&simon, Action{ActionType: Check})
	assert.NoError(t, err)
	assert.Equal(t, 40, simon.Stack)
	err = game.TakeAction(&chien, Action{ActionType: Bet, Value: 5})
	assert.NoError(t, err)
	assert.Equal(t, 85, chien.Stack)
	// Calls 5, bets 10. Total: 15
	err = game.TakeAction(&jarry, Action{ActionType: Bet, Value: 10})
	assert.NoError(t, err)
	assert.Equal(t, 5, jarry.Stack)
	err = game.TakeAction(&jason, Action{ActionType: Call})
	assert.NoError(t, err)
	assert.Equal(t, 0, jason.Stack)
	err = game.TakeAction(&simon, Action{ActionType: Call})
	assert.NoError(t, err)
	assert.Equal(t, 25, simon.Stack)
	err = game.TakeAction(&chien, Action{ActionType: Fold})
	assert.NoError(t, err)

	// Turn
	// Bet is 15 each, Jarry can only afford 5, so side pot becomes 20 (10 from flop + 5 (Jarry's all-in) * 2 (players))
	assert.Equal(t, Turn, game.GameState.Round)
	err = game.TakeAction(&simon, Action{ActionType: Bet, Value: 15})
	assert.NoError(t, err)
	assert.Equal(t, 10, simon.Stack)
	// Not in hand anymore
	err = game.TakeAction(&chien, Action{ActionType: Call})
	assert.Error(t, err)
	err = game.TakeAction(&jarry, Action{ActionType: Call})
	assert.NoError(t, err)
	assert.Equal(t, 0, jarry.Stack)
	// Jason is already all-in
	err = game.TakeAction(&jason, Action{ActionType: Call})
	assert.Error(t, err)

	// Simon bet 15, but only 5 was called, so 10 should have been added back to Simon's stack
	assert.Equal(t, 20, simon.Stack)

	// River
	// TODO: No more action, this should be auto done
	assert.Equal(t, River, game.GameState.Round)
	err = game.TakeAction(&simon, Action{ActionType: Check})
	assert.NoError(t, err)

	// HandEnd
	// First showdown is Simon vs Jarry. Jarry wins side pot of 20 and now has 20
	// Second showdown is Simon vs Jarry vs Jason. Jason wins main pot of 75 and now has 75
	assert.Equal(t, HandEnd, game.GameState.Round)
	assert.Equal(t, 75, jason.Stack)
	assert.Equal(t, 85, chien.Stack)
	assert.Equal(t, 20, simon.Stack)
	assert.Equal(t, 20, jarry.Stack)
}

func TestGamePreFlopOption(t *testing.T) {
	game := NewDeterministicGame(1, 2, getDeck)
	jason := NewPlayer("Jason", 100)
	err := game.SitPlayer(&jason, 2)
	assert.NoError(t, err)
	simon := NewPlayer("Simon", 100)
	err = game.SitPlayer(&simon, 5)
	assert.NoError(t, err)
	chien := NewPlayer("Chien", 100)
	err = game.SitPlayer(&chien, 7)
	assert.NoError(t, err)

	game.DealHand()

	assert.Equal(t, PreFlop, game.GameState.Round)
	err = game.TakeAction(&jason, Action{ActionType: Call})
	assert.NoError(t, err)
	err = game.TakeAction(&simon, Action{ActionType: Call})
	assert.NoError(t, err)
	// Option check
	err = game.TakeAction(&chien, Action{ActionType: Check})
	assert.NoError(t, err)

	assert.Equal(t, Flop, game.GameState.Round)
	err = game.TakeAction(&simon, Action{ActionType: Check})
	assert.NoError(t, err)
	assert.Equal(t, 98, jason.Stack)
	assert.Equal(t, 98, simon.Stack)
	assert.Equal(t, 98, chien.Stack)

	game = NewDeterministicGame(1, 2, getDeck)
	jason = NewPlayer("Jason", 100)
	err = game.SitPlayer(&jason, 2)
	assert.NoError(t, err)
	simon = NewPlayer("Simon", 100)
	err = game.SitPlayer(&simon, 5)
	assert.NoError(t, err)
	chien = NewPlayer("Chien", 100)
	err = game.SitPlayer(&chien, 7)
	assert.NoError(t, err)

	game.DealHand()

	assert.Equal(t, PreFlop, game.GameState.Round)
	err = game.TakeAction(&jason, Action{ActionType: Call})
	assert.NoError(t, err)
	err = game.TakeAction(&simon, Action{ActionType: Call})
	assert.NoError(t, err)
	// Option raise
	err = game.TakeAction(&chien, Action{ActionType: Bet, Value: 5})
	assert.NoError(t, err)
	// Still pre flop
	assert.Equal(t, PreFlop, game.GameState.Round)
	err = game.TakeAction(&jason, Action{ActionType: Call})
	assert.NoError(t, err)
	err = game.TakeAction(&simon, Action{ActionType: Call})
	assert.NoError(t, err)

	assert.Equal(t, Flop, game.GameState.Round)
	err = game.TakeAction(&simon, Action{ActionType: Check})
	assert.NoError(t, err)
	assert.Equal(t, 93, jason.Stack)
	assert.Equal(t, 93, simon.Stack)
	assert.Equal(t, 93, chien.Stack)
}

func TestGetPlayerState(t *testing.T) {
	game := NewDeterministicGame(1, 2, getDeck)
	jason := NewPlayer("Jason", 100)
	err := game.SitPlayer(&jason, 2)
	assert.NoError(t, err)
	simon := NewPlayer("Simon", 100)
	err = game.SitPlayer(&simon, 5)
	assert.NoError(t, err)
	chien := NewPlayer("Chien", 100)
	err = game.SitPlayer(&chien, 7)
	assert.NoError(t, err)

	game.DealHand()

	assert.Equal(t, PreFlop, game.GameState.Round)
	err = game.TakeAction(&jason, Action{ActionType: Call})
	assert.NoError(t, err)
	ps := game.GetPlayerState(&jason)
	fmt.Println(ps)
	err = game.TakeAction(&simon, Action{ActionType: Call})
	assert.NoError(t, err)
	// Option check
	err = game.TakeAction(&chien, Action{ActionType: Check})
	assert.NoError(t, err)

	assert.Equal(t, Flop, game.GameState.Round)
	err = game.TakeAction(&simon, Action{ActionType: Check})
	assert.NoError(t, err)
	assert.Equal(t, 98, jason.Stack)
	assert.Equal(t, 98, simon.Stack)
	assert.Equal(t, 98, chien.Stack)
}