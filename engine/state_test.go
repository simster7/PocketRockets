package engine

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func getNewTestGameState(seats []int) State {
	players := new(Players)
	for _, seat := range seats {
		players[seat] = &Player{Name: string(seat), Stack: 100}
	}
	return GetNewHandState(*players, 2, 2, 1, getDeck())
}

func TestProcessPots(t *testing.T) {

	// Test standard pot
	gs := getNewTestGameState([]int{2, 5, 7})
	gs.Players[2].Bet = 10
	gs.Players[5].Bet = 10
	gs.Players[7].Bet = 10

	gs.Pots = []int{0}
	gs.PotContenders = [][]int{AllPlayers}

	gs.processPots()

	assert.Equal(t, []int{30}, gs.Pots)
	assert.Equal(t, [][]int{{0, 1, 2, 3, 4, 5, 6, 7, 8}}, gs.PotContenders)

	// Test one all in
	gs = getNewTestGameState([]int{2, 5, 7})
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
	gs = getNewTestGameState([]int{0, 2, 5, 7})
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
	gs = getNewTestGameState([]int{0, 2, 3, 5, 7})
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
	gs = getNewTestGameState([]int{5, 7})
	gs.Players[5].Bet = 30
	gs.Players[7].Bet = 10
	gs.Players[7].IsAllIn = true

	gs.Pots = []int{0}
	gs.PotContenders = [][]int{AllPlayers}

	gs.processPots()

	assert.Equal(t, []int{20, 0}, gs.Pots)
	assert.Equal(t, [][]int{AllPlayers, {0, 1, 2, 3, 4, 5, 6, 8}}, gs.PotContenders)
}
