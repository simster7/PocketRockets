package engine

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func getNewTestGameState() GameState {
	jason := NewPlayer("Jason", 100)
	simon := NewPlayer("Simon", 50)
	chien := NewPlayer("Chien", 20)

	seats := [9]Seat{
		{Index: 0, Occupied: false, Player: nil},
		{Index: 1, Occupied: false, Player: nil},
		{Index: 2, Occupied: true, Player: &jason},
		{Index: 3, Occupied: false, Player: nil},
		{Index: 4, Occupied: false, Player: nil},
		{Index: 5, Occupied: true, Player: &simon},
		{Index: 6, Occupied: false, Player: nil},
		{Index: 7, Occupied: true, Player: &chien},
		{Index: 8, Occupied: false, Player: nil},
	}

	deck := getDeck()

	gs, _ := GetNewHandGameState(seats, 2, 2, 1, deck)
	return gs
}

func TestProcessPots(t *testing.T) {

	// Test standard pot
	betVector := [9]BetVectorNode{
		{Amount: 0, IsAllIn: false},
		{Amount: 0, IsAllIn: false},
		{Amount: 10, IsAllIn: false},
		{Amount: 0, IsAllIn: false},
		{Amount: 0, IsAllIn: false},
		{Amount: 10, IsAllIn: false},
		{Amount: 0, IsAllIn: false},
		{Amount: 10, IsAllIn: false},
		{Amount: 0, IsAllIn: false},
	}

	pots := []int{0}
	potContenders := [][]int{{2, 5, 7}}

	processPots(&betVector, &potContenders, &pots)

	assert.Equal(t, []int{30}, pots)
	assert.Equal(t, [][]int{{2, 5, 7}}, potContenders)
	assert.Equal(t, getZeroBetVector(), betVector)

	// Test one all in

	betVector = [9]BetVectorNode{
		{Amount: 0, IsAllIn: false},
		{Amount: 0, IsAllIn: false},
		{Amount: 30, IsAllIn: false},
		{Amount: 0, IsAllIn: false},
		{Amount: 0, IsAllIn: false},
		{Amount: 30, IsAllIn: false},
		{Amount: 0, IsAllIn: false},
		{Amount: 20, IsAllIn: true},
		{Amount: 0, IsAllIn: false},
	}

	pots = []int{0}
	potContenders =[][]int{{2, 5, 7}}

	processPots(&betVector, &potContenders, &pots)

	assert.Equal(t, []int{60, 20}, pots)
	assert.Equal(t, [][]int{{2, 5, 7}, {2, 5}}, potContenders)
	assert.Equal(t, getZeroBetVector(), betVector)


	// Test two all ins

	betVector = [9]BetVectorNode{
		{Amount: 40, IsAllIn: false},
		{Amount: 0, IsAllIn: false},
		{Amount: 40, IsAllIn: false},
		{Amount: 0, IsAllIn: false},
		{Amount: 0, IsAllIn: false},
		{Amount: 30, IsAllIn: true},
		{Amount: 0, IsAllIn: false},
		{Amount: 20, IsAllIn: true},
		{Amount: 0, IsAllIn: false},
	}

	pots = []int{0}
	potContenders =[][]int{{0, 2, 5, 7}}

	processPots(&betVector, &potContenders, &pots)

	assert.Equal(t, []int{80, 30, 20}, pots)
	assert.Equal(t, [][]int{{0, 2, 5, 7}, {0, 2, 5}, {0, 2}}, potContenders)
	assert.Equal(t, getZeroBetVector(), betVector)


	// Test two all ins with folds

	betVector = [9]BetVectorNode{
		{Amount: 40, IsAllIn: false},
		{Amount: 0, IsAllIn: false},
		{Amount: 40, IsAllIn: false},
		{Amount: 10, IsAllIn: false},
		{Amount: 0, IsAllIn: false},
		{Amount: 30, IsAllIn: true},
		{Amount: 0, IsAllIn: false},
		{Amount: 20, IsAllIn: true},
		{Amount: 0, IsAllIn: false},
	}

	pots = []int{0}
	potContenders =[][]int{{0, 2, 3, 5, 7}}

	processPots(&betVector, &potContenders, &pots)

	assert.Equal(t, []int{90, 30, 20}, pots)
	assert.Equal(t, [][]int{{0, 2, 3, 5, 7}, {0, 2, 3, 5}, {0, 2, 3}}, potContenders)
	assert.Equal(t, getZeroBetVector(), betVector)

}
