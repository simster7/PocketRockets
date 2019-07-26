package engine

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGameBasic (t *testing.T) {
	game := NewGame(1, 2)
	grace := NewPlayer("Grace", 100)
	err := game.SitPlayer(grace, 0)
	assert.NoError(t, err)
    jason := NewPlayer("Jason", 100)
	err = game.SitPlayer(jason, 1)
	assert.NoError(t, err)
    simon := NewPlayer("Simon", 100)
	err = game.SitPlayer(simon, 3)
	assert.NoError(t, err)
    hersh := NewPlayer("Hersh", 100)
	err = game.SitPlayer(hersh, 4)
	assert.NoError(t, err)
    chien := NewPlayer("Chien", 100)
	err = game.SitPlayer(chien, 5)
	assert.NoError(t, err)
    jarry := NewPlayer("Jarry", 100)
	err = game.SitPlayer(jarry, 6)
	assert.NoError(t, err)
}
