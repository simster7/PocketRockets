package engine

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGameBasic(t *testing.T) {
	// TODO Replace GameState calls with API calls
	game := NewDeterministicGame(1, 2)
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
	assert.Equal(t, 95, chien.Stack)
	err = game.TakeAction(&jarry, Action{ActionType: call})
	assert.NoError(t, err)
	assert.Equal(t, 95, jarry.Stack)
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
	assert.Equal(t, 85, jarry.Stack)
	// Can't check a bet
	err = game.TakeAction(&jason, Action{ActionType: check})
	assert.Error(t, err)
	err = game.TakeAction(&jason, Action{ActionType: fold})
	assert.NoError(t, err)
	assert.Equal(t, true, jason.Folded)
	err = game.TakeAction(&simon, Action{ActionType: call})
	assert.NoError(t, err)
	assert.Equal(t, 85, simon.Stack)
	err = game.TakeAction(&chien, Action{ActionType: call})
	assert.NoError(t, err)
	assert.Equal(t, 85, chien.Stack)

	// Turn
	assert.Equal(t, Turn, game.GameState.Round)
	assert.Len(t, game.GameState.getCommunityCards(), 4)

}
