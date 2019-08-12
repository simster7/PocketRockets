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
	assert.Equal(t, int32(PreFlop), game.GetPlayerState(&simon).BettingRound)
	assert.Nil(t, game.GameState.getCommunityCards())
	// Can't play out of turn
	err = game.TakeAction(&jarry, Action{ActionType: Bet, Value: 5})
	assert.Error(t, err)
	assert.Equal(t, int32(100), game.GetPlayerState(&jarry).Seats[jarry.SeatNumber].Player.Stack)
	err = game.TakeAction(&chien, Action{ActionType: Bet, Value: 5})
	assert.NoError(t, err)
	assert.Equal(t, int32(93), game.GetPlayerState(&chien).Seats[chien.SeatNumber].Player.Stack)
	err = game.TakeAction(&jarry, Action{ActionType: Call})
	assert.NoError(t, err)
	assert.Equal(t, int32(93), game.GetPlayerState(&jarry).Seats[jarry.SeatNumber].Player.Stack)
	err = game.TakeAction(&grace, Action{ActionType: Fold})
	assert.NoError(t, err)
	assert.Equal(t, int32(100), game.GetPlayerState(&grace).Seats[grace.SeatNumber].Player.Stack)
	assert.Equal(t, true, game.GetPlayerState(&grace).Seats[grace.SeatNumber].Player.Folded)
	err = game.TakeAction(&jason, Action{ActionType: Call})
	assert.NoError(t, err)
	err = game.TakeAction(&simon, Action{ActionType: Call})
	assert.NoError(t, err)
	err = game.TakeAction(&hersh, Action{ActionType: Fold})
	assert.NoError(t, err)
	assert.Equal(t, true, game.GetPlayerState(&hersh).Seats[hersh.SeatNumber].Player.Folded)

	// Flop
	assert.Equal(t, int32(Flop), game.GetPlayerState(&simon).BettingRound)
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
	assert.Equal(t, int32(83), game.GetPlayerState(&jarry).Seats[jarry.SeatNumber].Player.Stack)
	// Can't check a bet
	err = game.TakeAction(&jason, Action{ActionType: Check})
	assert.Error(t, err)
	err = game.TakeAction(&jason, Action{ActionType: Fold})
	assert.NoError(t, err)
	assert.Equal(t, true, game.GetPlayerState(&jason).Seats[jason.SeatNumber].Player.Folded)
	err = game.TakeAction(&simon, Action{ActionType: Call})
	assert.NoError(t, err)
	assert.Equal(t, int32(83), game.GetPlayerState(&simon).Seats[simon.SeatNumber].Player.Stack)
	err = game.TakeAction(&chien, Action{ActionType: Call})
	assert.NoError(t, err)
	assert.Equal(t, int32(83), game.GetPlayerState(&chien).Seats[chien.SeatNumber].Player.Stack)

	// Turn
	assert.Equal(t, int32(Turn), game.GetPlayerState(&simon).BettingRound)
	assert.Len(t, game.GameState.getCommunityCards(), 4)
	assert.Equal(t, game.GameState.Pots[0], 60)

	err = game.TakeAction(&simon, Action{ActionType: Check})
	assert.NoError(t, err)
	err = game.TakeAction(&chien, Action{ActionType: Check})
	assert.NoError(t, err)
	err = game.TakeAction(&jarry, Action{ActionType: Check})
	assert.NoError(t, err)

	// River
	assert.Equal(t, int32(River), game.GetPlayerState(&simon).BettingRound)
	assert.Len(t, game.GameState.getCommunityCards(), 5)
	assert.Equal(t, game.GameState.Pots[0], 60)

	err = game.TakeAction(&simon, Action{ActionType: Bet, Value: 10})
	assert.NoError(t, err)
	assert.Equal(t, int32(73), game.GetPlayerState(&simon).Seats[simon.SeatNumber].Player.Stack)
	err = game.TakeAction(&chien, Action{ActionType: Fold})
	assert.NoError(t, err)
	assert.Equal(t, true, game.GetPlayerState(&chien).Seats[chien.SeatNumber].Player.Folded)
	err = game.TakeAction(&jarry, Action{ActionType: Call})
	assert.NoError(t, err)

	// Post River
	// Board hits a flush, split pot
	assert.Equal(t, int32(HandEnd), game.GetPlayerState(&simon).BettingRound)
	assert.Equal(t, game.GameState.Pots[0], 80)
	assert.Len(t, game.GameState.getCommunityCards(), 5)
	assert.False(t, game.IsHandActive)
	assert.False(t, game.GameState.IsHandActive)
	assert.Equal(t, int32(113), game.GetPlayerState(&simon).Seats[simon.SeatNumber].Player.Stack)
	assert.Equal(t, int32(113), game.GetPlayerState(&jarry).Seats[jarry.SeatNumber].Player.Stack)
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
	assert.Equal(t, int32(PreFlop), game.GetPlayerState(&simon).BettingRound)
	err = game.TakeAction(&jason, Action{ActionType: Bet, Value: 8})
	assert.NoError(t, err)
	assert.Equal(t, int32(90), game.GetPlayerState(&jason).Seats[jason.SeatNumber].Player.Stack)
	// Can't check a call
	err = game.TakeAction(&simon, Action{ActionType: Check})
	assert.Error(t, err)
	err = game.TakeAction(&simon, Action{ActionType: Call})
	assert.NoError(t, err)
	assert.Equal(t, int32(90), game.GetPlayerState(&jason).Seats[jason.SeatNumber].Player.Stack)
	err = game.TakeAction(&chien, Action{ActionType: Call})
	assert.NoError(t, err)

	// Flop
	assert.Equal(t, int32(Flop), game.GetPlayerState(&simon).BettingRound)
	err = game.TakeAction(&simon, Action{ActionType: Check})
	assert.NoError(t, err)
	err = game.TakeAction(&chien, Action{ActionType: Check})
	assert.NoError(t, err)
	err = game.TakeAction(&jason, Action{ActionType: Check})
	assert.NoError(t, err)

	// Turn
	assert.Equal(t, int32(Turn), game.GetPlayerState(&simon).BettingRound)
	// Can't bet more than stack
	err = game.TakeAction(&simon, Action{ActionType: Bet, Value: 1000})
	assert.Error(t, err)
	err = game.TakeAction(&simon, Action{ActionType: Bet, Value: 10})
	assert.NoError(t, err)
	assert.Equal(t, int32(80), game.GetPlayerState(&simon).Seats[simon.SeatNumber].Player.Stack)
	err = game.TakeAction(&chien, Action{ActionType: Call})
	assert.NoError(t, err)
	assert.Equal(t, int32(80), game.GetPlayerState(&chien).Seats[chien.SeatNumber].Player.Stack)
	err = game.TakeAction(&jason, Action{ActionType: Fold})
	assert.NoError(t, err)

	// River
	assert.Equal(t, int32(River), game.GetPlayerState(&simon).BettingRound)
	err = game.TakeAction(&simon, Action{ActionType: Check})
	assert.NoError(t, err)
	err = game.TakeAction(&chien, Action{ActionType: Check})
	assert.NoError(t, err)

	assert.Equal(t, int32(HandEnd), game.GetPlayerState(&simon).BettingRound)
	assert.Equal(t, int32(80), game.GetPlayerState(&simon).Seats[simon.SeatNumber].Player.Stack)
	assert.Equal(t, int32(130), game.GetPlayerState(&chien).Seats[chien.SeatNumber].Player.Stack)

	assert.False(t, game.IsHandActive)

	jarry := NewPlayer("Jarry", 100)
	err = game.SitPlayer(&jarry, 8)
	assert.NoError(t, err)

	game.DealHand()

	// Pre flop
	assert.Equal(t, int32(PreFlop), game.GetPlayerState(&simon).BettingRound)
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
	assert.Equal(t, int32(Flop), game.GetPlayerState(&simon).BettingRound)
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
	assert.Equal(t, int32(Turn), game.GetPlayerState(&simon).BettingRound)
	assert.Len(t, game.GameState.getCommunityCards(), 4)

	err = game.TakeAction(&jason, Action{ActionType: Check})
	assert.NoError(t, err)
	err = game.TakeAction(&simon, Action{ActionType: Check})
	assert.NoError(t, err)

	// River
	assert.Equal(t, int32(River), game.GetPlayerState(&simon).BettingRound)
	assert.Len(t, game.GameState.getCommunityCards(), 5)

	err = game.TakeAction(&jason, Action{ActionType: Check})
	assert.NoError(t, err)
	err = game.TakeAction(&simon, Action{ActionType: Check})
	assert.NoError(t, err)

	// PostRiver
	assert.Equal(t, int32(HandEnd), game.GetPlayerState(&simon).BettingRound)

	assert.Equal(t, int32(115), game.GetPlayerState(&simon).Seats[simon.SeatNumber].Player.Stack)
	assert.Equal(t, int32(68), game.GetPlayerState(&jason).Seats[jason.SeatNumber].Player.Stack)
	assert.Equal(t, int32(88), game.GetPlayerState(&jarry).Seats[jarry.SeatNumber].Player.Stack)
	assert.Equal(t, int32(129), game.GetPlayerState(&chien).Seats[chien.SeatNumber].Player.Stack)

	game.DealHand()

	// Pre flop
	assert.Equal(t, int32(PreFlop), game.GetPlayerState(&simon).BettingRound)
	assert.Nil(t, game.GameState.getCommunityCards())

	err = game.TakeAction(&simon, Action{ActionType: Bet, Value: 10})
	assert.NoError(t, err)
	err = game.TakeAction(&chien, Action{ActionType: Call})
	assert.NoError(t, err)
	err = game.TakeAction(&jarry, Action{ActionType: Call})
	assert.NoError(t, err)
	err = game.TakeAction(&jason, Action{ActionType: Call})
	assert.NoError(t, err)

	assert.Equal(t, int32(103), game.GetPlayerState(&simon).Seats[simon.SeatNumber].Player.Stack)
	assert.Equal(t, int32(56), game.GetPlayerState(&jason).Seats[jason.SeatNumber].Player.Stack)
	assert.Equal(t, int32(76), game.GetPlayerState(&jarry).Seats[jarry.SeatNumber].Player.Stack)
	assert.Equal(t, int32(117), game.GetPlayerState(&chien).Seats[chien.SeatNumber].Player.Stack)

	// Pre flop
	assert.Equal(t, int32(Flop), game.GetPlayerState(&simon).BettingRound)
	// DONK!
	err = game.TakeAction(&jarry, Action{ActionType: Bet, Value: 50})
	assert.NoError(t, err)
	assert.Equal(t, int32(26), game.GetPlayerState(&jarry).Seats[jarry.SeatNumber].Player.Stack)
	err = game.TakeAction(&jason, Action{ActionType: Fold})
	assert.NoError(t, err)
	err = game.TakeAction(&simon, Action{ActionType: Fold})
	assert.NoError(t, err)
	err = game.TakeAction(&chien, Action{ActionType: Fold})
	assert.NoError(t, err)

	assert.False(t, game.IsHandActive)
	assert.Equal(t, int32(103), game.GetPlayerState(&simon).Seats[simon.SeatNumber].Player.Stack)
	assert.Equal(t, int32(56), game.GetPlayerState(&jason).Seats[jason.SeatNumber].Player.Stack)
	assert.Equal(t, int32(124), game.GetPlayerState(&jarry).Seats[jarry.SeatNumber].Player.Stack)
	assert.Equal(t, int32(117), game.GetPlayerState(&chien).Seats[chien.SeatNumber].Player.Stack)
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
	assert.Equal(t, int32(PreFlop), game.GetPlayerState(&simon).BettingRound)
	err = game.TakeAction(&jason, Action{ActionType: Bet, Value: 8})
	assert.NoError(t, err)
	assert.Equal(t, int32(10), game.GetPlayerState(&jason).Seats[jason.SeatNumber].Player.Stack)
	err = game.TakeAction(&simon, Action{ActionType: Call})
	assert.NoError(t, err)
	assert.Equal(t, int32(40), game.GetPlayerState(&simon).Seats[simon.SeatNumber].Player.Stack)
	err = game.TakeAction(&chien, Action{ActionType: Call})
	assert.NoError(t, err)
	assert.Equal(t, int32(90), game.GetPlayerState(&chien).Seats[chien.SeatNumber].Player.Stack)

	// Flop
	// Bet is 30 each, Jason can only afford 10. Main pot becomes 60 (30 from flop + 10 (Jason's all-in) * 3 (players))
	// Side pot becomes 40 (20 each from Simon and Chien)
	assert.Equal(t, int32(Flop), game.GetPlayerState(&simon).BettingRound)
	err = game.TakeAction(&simon, Action{ActionType: Bet, Value: 30})
	assert.NoError(t, err)
	assert.Equal(t, int32(10), game.GetPlayerState(&simon).Seats[simon.SeatNumber].Player.Stack)
	err = game.TakeAction(&chien, Action{ActionType: Call})
	assert.NoError(t, err)
	assert.Equal(t, int32(60), game.GetPlayerState(&chien).Seats[chien.SeatNumber].Player.Stack)
	err = game.TakeAction(&jason, Action{ActionType: Call})
	assert.NoError(t, err)
	assert.Equal(t, int32(0), game.GetPlayerState(&jason).Seats[jason.SeatNumber].Player.Stack)

	// Turn
	assert.Equal(t, int32(Turn), game.GetPlayerState(&simon).BettingRound)
	err = game.TakeAction(&simon, Action{ActionType: Check})
	assert.NoError(t, err)
	err = game.TakeAction(&chien, Action{ActionType: Check})
	assert.NoError(t, err)

	// River
	assert.Equal(t, int32(River), game.GetPlayerState(&simon).BettingRound)
	err = game.TakeAction(&simon, Action{ActionType: Check})
	assert.NoError(t, err)
	err = game.TakeAction(&chien, Action{ActionType: Check})
	assert.NoError(t, err)

	// HandEnd
	// Jason wins main pot of 60, Chien wins side pot of 40
	assert.Equal(t, int32(HandEnd), game.GetPlayerState(&simon).BettingRound)
	assert.Equal(t, int32(60), game.GetPlayerState(&jason).Seats[jason.SeatNumber].Player.Stack)
	assert.Equal(t, int32(100), game.GetPlayerState(&chien).Seats[chien.SeatNumber].Player.Stack)
	assert.Equal(t, int32(10), game.GetPlayerState(&simon).Seats[simon.SeatNumber].Player.Stack)
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
	assert.Equal(t, int32(PreFlop), game.GetPlayerState(&simon).BettingRound)
	err = game.TakeAction(&jarry, Action{ActionType: Bet, Value: 8})
	assert.NoError(t, err)
	assert.Equal(t, int32(20), game.GetPlayerState(&jason).Seats[jason.SeatNumber].Player.Stack)
	err = game.TakeAction(&jason, Action{ActionType: Call})
	assert.NoError(t, err)
	assert.Equal(t, int32(10), game.GetPlayerState(&jason).Seats[jason.SeatNumber].Player.Stack)
	err = game.TakeAction(&simon, Action{ActionType: Call})
	assert.NoError(t, err)
	assert.Equal(t, int32(40), game.GetPlayerState(&simon).Seats[simon.SeatNumber].Player.Stack)
	err = game.TakeAction(&chien, Action{ActionType: Call})
	assert.NoError(t, err)
	assert.Equal(t, int32(90), game.GetPlayerState(&chien).Seats[chien.SeatNumber].Player.Stack)

	// Flop
	// Bet is 15 each, Jason can only afford 10. Main pot becomes 80 (40 from preflop + 10 (Jason's all-in) * 4 (players))
	// First side pot becomes 15 (5 each from Simon, Chien, and Jarry)
	assert.Equal(t, int32(Flop), game.GetPlayerState(&simon).BettingRound)
	err = game.TakeAction(&simon, Action{ActionType: Bet, Value: 15})
	assert.NoError(t, err)
	assert.Equal(t, int32(25), game.GetPlayerState(&simon).Seats[simon.SeatNumber].Player.Stack)
	err = game.TakeAction(&chien, Action{ActionType: Call})
	assert.NoError(t, err)
	assert.Equal(t, int32(75), game.GetPlayerState(&chien).Seats[chien.SeatNumber].Player.Stack)
	err = game.TakeAction(&jarry, Action{ActionType: Call})
	assert.NoError(t, err)
	assert.Equal(t, int32(5), game.GetPlayerState(&jarry).Seats[jarry.SeatNumber].Player.Stack)
	err = game.TakeAction(&jason, Action{ActionType: Call})
	assert.NoError(t, err)
	assert.Equal(t, int32(0), game.GetPlayerState(&jason).Seats[jason.SeatNumber].Player.Stack)

	// Turn
	// Bet is 15 each, Jarry can only afford 5. First side pot becomes 30 (15 from flop + 5 (Jarry's all-in) * 3 (players))
	// Second side pot becomes 20 (10 each from Simon and Chien)
	assert.Equal(t, int32(Turn), game.GetPlayerState(&simon).BettingRound)
	err = game.TakeAction(&simon, Action{ActionType: Bet, Value: 15})
	assert.NoError(t, err)
	assert.Equal(t, int32(10), game.GetPlayerState(&simon).Seats[simon.SeatNumber].Player.Stack)
	err = game.TakeAction(&chien, Action{ActionType: Call})
	assert.NoError(t, err)
	assert.Equal(t, int32(60), game.GetPlayerState(&chien).Seats[chien.SeatNumber].Player.Stack)
	err = game.TakeAction(&jarry, Action{ActionType: Call})
	assert.NoError(t, err)
	assert.Equal(t, int32(0), game.GetPlayerState(&jarry).Seats[jarry.SeatNumber].Player.Stack)
	// Jason is already all-in
	err = game.TakeAction(&jason, Action{ActionType: Call})
	assert.Error(t, err)

	// River
	assert.Equal(t, int32(River), game.GetPlayerState(&simon).BettingRound)
	err = game.TakeAction(&simon, Action{ActionType: Check})
	assert.NoError(t, err)
	err = game.TakeAction(&chien, Action{ActionType: Check})
	assert.NoError(t, err)

	// HandEnd
	// First showdown is Simon vs Chien. Chien wins second side pot of 20 and now has 80
	// Second showdown is Simon vs Chien vs Jarry. Jarry wins first sidepot of 30 and now has 30
	// Last showdown is family pot. Jason wins main pot of 80 and now has 80
	assert.Equal(t, int32(HandEnd), game.GetPlayerState(&simon).BettingRound)
	assert.Equal(t, int32(80), game.GetPlayerState(&jason).Seats[jason.SeatNumber].Player.Stack)
	assert.Equal(t, int32(80), game.GetPlayerState(&chien).Seats[chien.SeatNumber].Player.Stack)
	assert.Equal(t, int32(10), game.GetPlayerState(&simon).Seats[simon.SeatNumber].Player.Stack)
	assert.Equal(t, int32(30), game.GetPlayerState(&jarry).Seats[jarry.SeatNumber].Player.Stack)

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
	assert.Equal(t, int32(PreFlop), game.GetPlayerState(&simon).BettingRound)
	err = game.TakeAction(&jarry, Action{ActionType: Bet, Value: 8})
	assert.NoError(t, err)
	assert.Equal(t, int32(20), game.GetPlayerState(&jason).Seats[jason.SeatNumber].Player.Stack)
	err = game.TakeAction(&jason, Action{ActionType: Call})
	assert.NoError(t, err)
	assert.Equal(t, int32(10), game.GetPlayerState(&jason).Seats[jason.SeatNumber].Player.Stack)
	err = game.TakeAction(&simon, Action{ActionType: Call})
	assert.NoError(t, err)
	assert.Equal(t, int32(40), game.GetPlayerState(&simon).Seats[simon.SeatNumber].Player.Stack)
	err = game.TakeAction(&chien, Action{ActionType: Call})
	assert.NoError(t, err)
	assert.Equal(t, int32(90), game.GetPlayerState(&chien).Seats[chien.SeatNumber].Player.Stack)

	// Flop
	// Main pot becomes 75 (40 from preflop + 10 (Jason's all-in) * 3 players + 5 (chien's bet-fold))
	// First side pot becomes 10 (5 each from simon, jarry)
	assert.Equal(t, int32(Flop), game.GetPlayerState(&simon).BettingRound)
	err = game.TakeAction(&simon, Action{ActionType: Check})
	assert.NoError(t, err)
	assert.Equal(t, int32(40), game.GetPlayerState(&simon).Seats[simon.SeatNumber].Player.Stack)
	err = game.TakeAction(&chien, Action{ActionType: Bet, Value: 5})
	assert.NoError(t, err)
	assert.Equal(t, int32(85), game.GetPlayerState(&chien).Seats[chien.SeatNumber].Player.Stack)
	// Calls 5, bets 10. Total: 15
	err = game.TakeAction(&jarry, Action{ActionType: Bet, Value: 10})
	assert.NoError(t, err)
	assert.Equal(t, int32(5), game.GetPlayerState(&jarry).Seats[jarry.SeatNumber].Player.Stack)
	err = game.TakeAction(&jason, Action{ActionType: Call})
	assert.NoError(t, err)
	assert.Equal(t, int32(0), game.GetPlayerState(&jason).Seats[jason.SeatNumber].Player.Stack)
	err = game.TakeAction(&simon, Action{ActionType: Call})
	assert.NoError(t, err)
	assert.Equal(t, int32(25), game.GetPlayerState(&simon).Seats[simon.SeatNumber].Player.Stack)
	err = game.TakeAction(&chien, Action{ActionType: Fold})
	assert.NoError(t, err)

	// Turn
	// Bet is 15 each, Jarry can only afford 5, so side pot becomes 20 (10 from flop + 5 (Jarry's all-in) * 2 (players))
	assert.Equal(t, int32(Turn), game.GetPlayerState(&simon).BettingRound)
	err = game.TakeAction(&simon, Action{ActionType: Bet, Value: 15})
	assert.NoError(t, err)
	assert.Equal(t, int32(10), game.GetPlayerState(&simon).Seats[simon.SeatNumber].Player.Stack)
	// Not in hand anymore
	err = game.TakeAction(&chien, Action{ActionType: Call})
	assert.Error(t, err)
	err = game.TakeAction(&jarry, Action{ActionType: Call})
	assert.NoError(t, err)
	assert.Equal(t, int32(0), game.GetPlayerState(&jarry).Seats[jarry.SeatNumber].Player.Stack)
	// Jason is already all-in
	err = game.TakeAction(&jason, Action{ActionType: Call})
	assert.Error(t, err)

	// Simon bet 15, but only 5 was called, so 10 should have been added back to Simon's stack
	assert.Equal(t, int32(20), game.GetPlayerState(&simon).Seats[simon.SeatNumber].Player.Stack)

	// River
	// TODO: No more action, this should be auto done
	assert.Equal(t, int32(River), game.GetPlayerState(&simon).BettingRound)
	err = game.TakeAction(&simon, Action{ActionType: Check})
	assert.NoError(t, err)

	// HandEnd
	// First showdown is Simon vs Jarry. Jarry wins side pot of 20 and now has 20
	// Second showdown is Simon vs Jarry vs Jason. Jason wins main pot of 75 and now has 75
	assert.Equal(t, int32(HandEnd), game.GetPlayerState(&simon).BettingRound)
	assert.Equal(t, int32(75), game.GetPlayerState(&jason).Seats[jason.SeatNumber].Player.Stack)
	assert.Equal(t, int32(85), game.GetPlayerState(&chien).Seats[chien.SeatNumber].Player.Stack)
	assert.Equal(t, int32(20), game.GetPlayerState(&simon).Seats[simon.SeatNumber].Player.Stack)
	assert.Equal(t, int32(20), game.GetPlayerState(&jarry).Seats[jarry.SeatNumber].Player.Stack)
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

	assert.Equal(t, int32(PreFlop), game.GetPlayerState(&simon).BettingRound)
	err = game.TakeAction(&jason, Action{ActionType: Call})
	assert.NoError(t, err)
	err = game.TakeAction(&simon, Action{ActionType: Call})
	assert.NoError(t, err)
	// Option check
	err = game.TakeAction(&chien, Action{ActionType: Check})
	assert.NoError(t, err)

	assert.Equal(t, int32(Flop), game.GetPlayerState(&simon).BettingRound)
	err = game.TakeAction(&simon, Action{ActionType: Check})
	assert.NoError(t, err)
	assert.Equal(t, int32(98), game.GetPlayerState(&jason).Seats[jason.SeatNumber].Player.Stack)
	assert.Equal(t, int32(98), game.GetPlayerState(&simon).Seats[simon.SeatNumber].Player.Stack)
	assert.Equal(t, int32(98), game.GetPlayerState(&chien).Seats[chien.SeatNumber].Player.Stack)

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

	assert.Equal(t, int32(PreFlop), game.GetPlayerState(&simon).BettingRound)
	err = game.TakeAction(&jason, Action{ActionType: Call})
	assert.NoError(t, err)
	err = game.TakeAction(&simon, Action{ActionType: Call})
	assert.NoError(t, err)
	// Option raise
	err = game.TakeAction(&chien, Action{ActionType: Bet, Value: 5})
	assert.NoError(t, err)
	// Still pre flop
	assert.Equal(t, int32(PreFlop), game.GetPlayerState(&simon).BettingRound)
	err = game.TakeAction(&jason, Action{ActionType: Call})
	assert.NoError(t, err)
	err = game.TakeAction(&simon, Action{ActionType: Call})
	assert.NoError(t, err)

	assert.Equal(t, int32(Flop), game.GetPlayerState(&simon).BettingRound)
	err = game.TakeAction(&simon, Action{ActionType: Check})
	assert.NoError(t, err)
	assert.Equal(t, int32(93), game.GetPlayerState(&jason).Seats[jason.SeatNumber].Player.Stack)
	assert.Equal(t, int32(93), game.GetPlayerState(&simon).Seats[simon.SeatNumber].Player.Stack)
	assert.Equal(t, int32(93), game.GetPlayerState(&chien).Seats[chien.SeatNumber].Player.Stack)
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

	assert.Equal(t, int32(PreFlop), game.GetPlayerState(&simon).BettingRound)
	err = game.TakeAction(&jason, Action{ActionType: Call})
	assert.NoError(t, err)
	ps := game.GetPlayerState(&jason)
	fmt.Println(ps)
	err = game.TakeAction(&simon, Action{ActionType: Call})
	assert.NoError(t, err)
	// Option check
	err = game.TakeAction(&chien, Action{ActionType: Check})
	assert.NoError(t, err)

	assert.Equal(t, int32(Flop), game.GetPlayerState(&simon).BettingRound)
	err = game.TakeAction(&simon, Action{ActionType: Check})
	assert.NoError(t, err)
	assert.Equal(t, int32(98), game.GetPlayerState(&jason).Seats[jason.SeatNumber].Player.Stack)
	assert.Equal(t, int32(98), game.GetPlayerState(&simon).Seats[simon.SeatNumber].Player.Stack)
	assert.Equal(t, int32(98), game.GetPlayerState(&chien).Seats[chien.SeatNumber].Player.Stack)
}
