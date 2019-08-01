package engine

import (
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
	gs := getNewTestGameState()

	gs.BetVector = [9]BetVectorNode{
		
	}
}
