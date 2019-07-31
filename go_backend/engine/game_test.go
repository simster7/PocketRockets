package engine

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGameBasic(t *testing.T) {
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
	err = game.TakeAction(&jarry, Action{ActionType: bet, Value: 5})
	assert.Error(t, err)
	assert.Equal(t, 100, jarry.Stack)
	err = game.TakeAction(&chien, Action{ActionType: bet, Value: 5})
	assert.NoError(t, err)
	assert.Equal(t, 93, chien.Stack)
	err = game.TakeAction(&jarry, Action{ActionType: call})
	assert.NoError(t, err)
	assert.Equal(t, 93, jarry.Stack)
	err = game.TakeAction(&grace, Action{ActionType: fold})
	assert.NoError(t, err)
	assert.Equal(t, 100, grace.Stack)
	assert.Equal(t, true, grace.Folded)
	err = game.TakeAction(&jason, Action{ActionType: call})
	assert.NoError(t, err)
	err = game.TakeAction(&simon, Action{ActionType: call})
	assert.NoError(t, err)
	err = game.TakeAction(&hersh, Action{ActionType: fold})
	assert.NoError(t, err)
	assert.Equal(t, true, hersh.Folded)

	// Flop
	assert.Equal(t, Flop, game.GameState.Round)
	assert.NotNil(t, game.GameState.getCommunityCards())
	assert.Len(t, game.GameState.getCommunityCards(), 3)
	assert.Equal(t, game.GameState.Pots[0], 30)

	// Can't call at start of round
	err = game.TakeAction(&simon, Action{ActionType: call})
	assert.Error(t, err)
	// Can't play when folded
	err = game.TakeAction(&hersh, Action{ActionType: call})
	assert.Error(t, err)
	err = game.TakeAction(&simon, Action{ActionType: check})
	assert.NoError(t, err)
	err = game.TakeAction(&chien, Action{ActionType: check})
	assert.NoError(t, err)
	err = game.TakeAction(&jarry, Action{ActionType: bet, Value: 10})
	assert.NoError(t, err)
	assert.Equal(t, 83, jarry.Stack)
	// Can't check a bet
	err = game.TakeAction(&jason, Action{ActionType: check})
	assert.Error(t, err)
	err = game.TakeAction(&jason, Action{ActionType: fold})
	assert.NoError(t, err)
	assert.Equal(t, true, jason.Folded)
	err = game.TakeAction(&simon, Action{ActionType: call})
	assert.NoError(t, err)
	assert.Equal(t, 83, simon.Stack)
	err = game.TakeAction(&chien, Action{ActionType: call})
	assert.NoError(t, err)
	assert.Equal(t, 83, chien.Stack)

	// Turn
	assert.Equal(t, Turn, game.GameState.Round)
	assert.Len(t, game.GameState.getCommunityCards(), 4)
	assert.Equal(t, game.GameState.Pots[0], 60)

	err = game.TakeAction(&simon, Action{ActionType: check})
	assert.NoError(t, err)
	err = game.TakeAction(&chien, Action{ActionType: check})
	assert.NoError(t, err)
	err = game.TakeAction(&jarry, Action{ActionType: check})
	assert.NoError(t, err)

	// River
	assert.Equal(t, River, game.GameState.Round)
	assert.Len(t, game.GameState.getCommunityCards(), 5)
	assert.Equal(t, game.GameState.Pots[0], 60)

	err = game.TakeAction(&simon, Action{ActionType: bet, Value: 10})
	assert.NoError(t, err)
	assert.Equal(t, 73, simon.Stack)
	err = game.TakeAction(&chien, Action{ActionType: fold})
	assert.NoError(t, err)
	assert.Equal(t, true, chien.Folded)
	err = game.TakeAction(&jarry, Action{ActionType: call})
	assert.NoError(t, err)
	assert.Equal(t, 73, jarry.Stack)

	// Post River
	assert.Equal(t, HandEnd, game.GameState.Round)
	assert.Equal(t, game.GameState.Pots[0], 80)
	assert.Len(t, game.GameState.getCommunityCards(), 5)
	assert.False(t, game.IsHandActive)
	assert.False(t, game.GameState.IsHandActive)
	assert.Equal(t, 153, simon.Stack)
	assert.Equal(t, 73, jarry.Stack)
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
	err = game.TakeAction(&jason, Action{ActionType: bet, Value: 8})
	assert.NoError(t, err)
	assert.Equal(t, 90, jason.Stack)
	// Can't check a call
	err = game.TakeAction(&simon, Action{ActionType: check})
	assert.Error(t, err)
	err = game.TakeAction(&simon, Action{ActionType: call})
	assert.NoError(t, err)
	assert.Equal(t, 90, jason.Stack)
	err = game.TakeAction(&chien, Action{ActionType: call})
	assert.NoError(t, err)

	// Flop
	assert.Equal(t, Flop, game.GameState.Round)
	err = game.TakeAction(&simon, Action{ActionType: check})
	assert.NoError(t, err)
	err = game.TakeAction(&chien, Action{ActionType: check})
	assert.NoError(t, err)
	err = game.TakeAction(&jason, Action{ActionType: check})
	assert.NoError(t, err)

	// Turn
	assert.Equal(t, Turn, game.GameState.Round)
	// Can't bet more than stack
	err = game.TakeAction(&simon, Action{ActionType: bet, Value: 1000})
	assert.Error(t, err)
	err = game.TakeAction(&simon, Action{ActionType: bet, Value: 10})
	assert.NoError(t, err)
	assert.Equal(t, 80, simon.Stack)
	err = game.TakeAction(&chien, Action{ActionType: call})
	assert.NoError(t, err)
	assert.Equal(t, 80, chien.Stack)
	err = game.TakeAction(&jason, Action{ActionType: fold})
	assert.NoError(t, err)

	// River
	assert.Equal(t, River, game.GameState.Round)
	err = game.TakeAction(&simon, Action{ActionType: check})
	assert.NoError(t, err)
	err = game.TakeAction(&chien, Action{ActionType: check})
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
	err = game.TakeAction(&simon, Action{ActionType: bet, Value: 5})
	assert.Error(t, err)
	err = game.TakeAction(&chien, Action{ActionType: bet, Value: 5})
	assert.Error(t, err)
	err = game.TakeAction(&jarry, Action{ActionType: bet, Value: 5})
	assert.Error(t, err)

	err = game.TakeAction(&jason, Action{ActionType: call})
	assert.NoError(t, err)
	err = game.TakeAction(&simon, Action{ActionType: bet, Value: 10})
	assert.NoError(t, err)
	err = game.TakeAction(&chien, Action{ActionType: fold})
	assert.NoError(t, err)
	err = game.TakeAction(&jarry, Action{ActionType: call})
	assert.NoError(t, err)
	err = game.TakeAction(&jason, Action{ActionType: call})
	assert.NoError(t, err)


	// Flop 37
	assert.Equal(t, Flop, game.GameState.Round)
	assert.Len(t, game.GameState.getCommunityCards(), 3)

	// Can't play out of turn
	err = game.TakeAction(&simon, Action{ActionType: bet, Value: 5})
	assert.Error(t, err)
	err = game.TakeAction(&jason, Action{ActionType: bet, Value: 5})
	assert.Error(t, err)

	err = game.TakeAction(&jarry, Action{ActionType: check})
	assert.NoError(t, err)
	err = game.TakeAction(&jason, Action{ActionType: check})
	assert.NoError(t, err)
	err = game.TakeAction(&simon, Action{ActionType: bet, Value: 10})
	assert.NoError(t, err)
	err = game.TakeAction(&jarry, Action{ActionType: fold})
	assert.NoError(t, err)
	assert.True(t, jarry.Folded)
	err = game.TakeAction(&jason, Action{ActionType: call})
	assert.NoError(t, err)

	// Turn 57
	assert.Equal(t, Turn, game.GameState.Round)
	assert.Len(t, game.GameState.getCommunityCards(), 4)

	err = game.TakeAction(&jason, Action{ActionType: check})
	assert.NoError(t, err)
	err = game.TakeAction(&simon, Action{ActionType: check})
	assert.NoError(t, err)

	// River
	assert.Equal(t, River, game.GameState.Round)
	assert.Len(t, game.GameState.getCommunityCards(), 5)

	err = game.TakeAction(&jason, Action{ActionType: check})
	assert.NoError(t, err)
	err = game.TakeAction(&simon, Action{ActionType: check})
	assert.NoError(t, err)


	// PostRiver
	assert.Equal(t, HandEnd, game.GameState.Round)

	assert.Equal(t, 115, simon.Stack)
	assert.Equal(t, 68, jason.Stack)
	assert.Equal(t, 88, jarry.Stack)
	assert.Equal(t, 129, chien.Stack)
}

func TestGameAllInSimple(t *testing.T) {
	game := NewDeterministicGame(1, 2, getDeck)
	jason := NewPlayer("Jason", 100)
	err := game.SitPlayer(&jason, 2)
	assert.NoError(t, err)
	simon := NewPlayer("Simon", 50)
	err = game.SitPlayer(&simon, 5)
	assert.NoError(t, err)
	chien := NewPlayer("Chien", 20)
	err = game.SitPlayer(&chien, 7)
	assert.NoError(t, err)

	game.DealHand()

	// Pre flop
	assert.Equal(t, PreFlop, game.GameState.Round)
	err = game.TakeAction(&jason, Action{ActionType: bet, Value: 8})
	assert.NoError(t, err)
	assert.Equal(t, 90, jason.Stack)
	err = game.TakeAction(&simon, Action{ActionType: call})
	assert.NoError(t, err)
	assert.Equal(t, 40, simon.Stack)
	err = game.TakeAction(&chien, Action{ActionType: call})
	assert.NoError(t, err)
	assert.Equal(t, 10, chien.Stack)

	// Flop
	assert.Equal(t, Flop, game.GameState.Round)
	err = game.TakeAction(&simon, Action{ActionType: bet, Value: 40})
	assert.NoError(t, err)
	assert.Equal(t, 0, simon.Stack)
	err = game.TakeAction(&chien, Action{ActionType: call})
	assert.NoError(t, err)
	assert.Equal(t, 0, chien.Stack)
	err = game.TakeAction(&jason, Action{ActionType: call})
	assert.NoError(t, err)
	assert.Equal(t, 50, jason.Stack)

}
