package engine

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetRankFrequencies(t *testing.T) {
	hand := GenerateHand("8S 7H 7S 7D 3S 3C 2C 2S")
	handRanks := getCardRankIdSlice(hand)
	pairs := getFrequencies(handRanks, 2, 2)
	assert.Equal(t, []int{1, 0}, pairs)

	hand = GenerateHand("8S 7H 7S 7D 3S KC 9C 2S")
	handRanks = getCardRankIdSlice(hand)
	pairs = getFrequencies(handRanks, 2, 2)
	assert.Nil(t, pairs)

	hand = GenerateHand("8S 7H 7S 7D 3S 3C 2C 2S")
	handRanks = getCardRankIdSlice(hand)
	trips := getFrequencies(handRanks, 3, 3)
	assert.Equal(t, []int{5}, trips)

	hand = GenerateHand("7C 7H 7S 7D 3S 3C 2C 2S")
	handRanks = getCardRankIdSlice(hand)
	fours := getFrequencies(handRanks, 4, 4)
	assert.Equal(t, []int{5}, fours)
}

func TestCheckHighCard(t *testing.T) {
	// High card evaluated correctly
	hand := GenerateHand("8S 7H TH KD 4C")
	match, tiebreakers := CheckHighCard(hand)
	assert.True(t, match)
	assert.Equal(t, Tiebreakers{11, 8, 6, 5, 2}, tiebreakers)

	// High card evaluated correctly even for seven cards
	hand = GenerateHand("8S 7H TH KD 4C 3C 2C")
	match, tiebreakers = CheckHighCard(hand)
	assert.True(t, match)
	assert.Equal(t, Tiebreakers{11, 8, 6, 5, 2}, tiebreakers)
}

func TestCheckPair(t *testing.T) {
	// One pair correctly evaluated
	hand := GenerateHand("8S 8H TH KD 4C")
	match, result := CheckPair(hand)
	assert.True(t, match)
	assert.Equal(t, Tiebreakers{6, 11, 8, 2}, result)

	// One pair correctly evaluated with seven cards
	hand = GenerateHand("8S 8H KH JD 4C AH 9C")
	match, result = CheckPair(hand)
	assert.True(t, match)
	assert.Equal(t, Tiebreakers{6, 12, 11, 9}, result)

	// No pair not evaluated, even with seven cards
	hand = GenerateHand("8S 7H KH TD QC 2S 4H")
	match, result = CheckPair(hand)
	assert.False(t, match)
	assert.Nil(t, result)
}

func TestCheckTwoPair(t *testing.T) {
	// Two pair correctly evaluated
	hand := GenerateHand("8S 8H TH TD 4C")
	match, result := CheckTwoPair(hand)
	assert.True(t, match)
	assert.Equal(t, Tiebreakers{8, 6, 2}, result)

	// Two pair correctly evaluated with seven cards
	hand = GenerateHand("8S 8H TH TD 4C 6C 9S")
	match, result = CheckTwoPair(hand)
	assert.True(t, match)
	assert.Equal(t, Tiebreakers{8, 6, 7}, result)

	// Two pair correctly evaluated with seven cards and presence of three pair
	hand = GenerateHand("8S 8H TH TD 9C 9S 6C")
	match, result = CheckTwoPair(hand)
	assert.True(t, match)
	assert.Equal(t, Tiebreakers{8, 7, 6}, result)

	// No two pair evaluated, even with seven cards
	hand = GenerateHand("8S 7H KH TD QC 2S 4H")
	match, result = CheckTwoPair(hand)
	assert.False(t, match)
	assert.Nil(t, result)

	// No two pair evaluated, even with seven cards and a single pair
	hand = GenerateHand("8S 8H KH TD QC 2S 4H")
	match, result = CheckTwoPair(hand)
	assert.False(t, match)
	assert.Nil(t, result)

	hand = GenerateHand("AC JS AS QH QC 7C 8D")
	match, result = CheckTwoPair(hand)
	assert.True(t, match)
	assert.Equal(t, Tiebreakers{12, 10, 9}, result)
}

func TestCheckThreeOfAKind(t *testing.T) {
	// Three of a kind correctly evaluated
	hand := GenerateHand("8S 8H 8C TD 4C")
	match, result := CheckThreeOfAKind(hand)
	assert.True(t, match)
	assert.Equal(t, Tiebreakers{6, 8, 2}, result)

	// Three of a kind correctly evaluated, even with seven cards
	hand = GenerateHand("8S 8H 8C TD 4C JC 2S")
	match, result = CheckThreeOfAKind(hand)
	assert.True(t, match)
	assert.Equal(t, Tiebreakers{6, 9, 8}, result)

	// No three of a kind evaluated, even with seven cards
	hand = GenerateHand("8S 8H 7C TD 4C JC 2S")
	match, result = CheckThreeOfAKind(hand)
	assert.False(t, match)
	assert.Nil(t, result)

	// Three of a kind correctly evaluated, even with seven cards and extra pair (full house)
	// This tests undefined behavior, so if it fails it is safe to delete
	hand = GenerateHand("8S 8H 8C TD TC JC 2S")
	match, result = CheckThreeOfAKind(hand)
	assert.True(t, match)
	assert.Equal(t, Tiebreakers{6, 9, 8}, result)
}

func TestCheckStraight(t *testing.T) {
	// Straight check works correctly
	hand := GenerateHand("2S 3H 4C 5D 6C")
	match, result := CheckStraight(hand)
	assert.True(t, match)
	assert.Equal(t, Tiebreakers{4}, result)

	// Straight check works correctly, even with seven cards
	hand = GenerateHand("2S 3H 4C 5D 6C KC JS")
	match, result = CheckStraight(hand)
	assert.True(t, match)
	assert.Equal(t, Tiebreakers{4}, result)

	// traight check works correctly, even with seven cards and straight is more than five cards
	hand = GenerateHand("2S 3H 4C 5D 6C 7C JS")
	match, result = CheckStraight(hand)
	assert.True(t, match)
	assert.Equal(t, Tiebreakers{5}, result)

	// Straight check works correctly, even with seven cards and is ace low
	hand = GenerateHand("2S 3H 4C 5D AC 9C JS")
	match, result = CheckStraight(hand)
	assert.True(t, match)
	assert.Equal(t, Tiebreakers{3}, result)

	// Straight check works correctly, even with seven cards and is ace high
	hand = GenerateHand("TS QH KC 5D AC 9C JS")
	match, result = CheckStraight(hand)
	assert.True(t, match)
	assert.Equal(t, Tiebreakers{12}, result)

	// No straight when there is not straight
	hand = GenerateHand("8S QH KC 6D AC 9C JS")
	match, result = CheckStraight(hand)
	assert.False(t, match)
	assert.Nil(t, result)

	// No straight wrap-around
	hand = GenerateHand("8S QH KC 6D AC 2C JS")
	match, result = CheckStraight(hand)
	assert.False(t, match)
	assert.Nil(t, result)

	hand = GenerateHand("AC JS AS QH QC 7C 8D")
	match, result = CheckStraight(hand)
	assert.False(t, match)
	assert.Nil(t, result)
}

func TestCheckFlush(t *testing.T) {
	// Flush check works correctly
	hand := GenerateHand("JS 3S TS 5S 6S")
	match, result := CheckFlush(hand)
	assert.True(t, match)
	assert.Equal(t, Tiebreakers{9, 8, 4, 3, 1}, result)

	// Flush check works correctly with seven cards, gets high card correctly
	hand = GenerateHand("JS 3S TS 5S 6S AS 8C")
	match, result = CheckFlush(hand)
	assert.True(t, match)
	assert.Equal(t, Tiebreakers{12, 9, 8, 4, 3}, result)

	// Flush not detected when there is no flush
	hand = GenerateHand("JS 3S TS 5S 6C AC 8C")
	match, result = CheckFlush(hand)
	assert.False(t, match)
	assert.Nil(t, result)
}

func TestCheckFullHouse(t *testing.T) {
	// Full house check works correctly
	hand := GenerateHand("JS JC JD 6C 6S")
	match, result := CheckFullHouse(hand)
	assert.True(t, match)
	assert.Equal(t, Tiebreakers{9, 4}, result)

	// Full house check works correctly with seven cards
	hand = GenerateHand("AS AC AD 6C 6S 2S KC")
	match, result = CheckFullHouse(hand)
	assert.True(t, match)
	assert.Equal(t, Tiebreakers{12, 4}, result)

	// Full house check works correctly with seven cards and two three of a kinds
	hand = GenerateHand("JS JC JD 6C 6S 6D KC")
	match, result = CheckFullHouse(hand)
	assert.True(t, match)
	assert.Equal(t, Tiebreakers{9, 4}, result)

	// Full house check works correctly with seven cards and two three of a kinds, order agnostic
	hand = GenerateHand("6C 6S 6D JS JC JD KC")
	match, result = CheckFullHouse(hand)
	assert.True(t, match)
	assert.Equal(t, Tiebreakers{9, 4}, result)

	// No full house when there is no full house
	hand = GenerateHand("JS JC KD 6C 6S 5D KC")
	match, result = CheckFullHouse(hand)
	assert.False(t, match)
	assert.Nil(t, result)
}

func TestCheckFourOfAKind(t *testing.T) {
	// Four of a kind check works correctly
	hand := GenerateHand("JS JC JD JH AS")
	match, result := CheckFourOfAKind(hand)
	assert.True(t, match)
	assert.Equal(t, Tiebreakers{9}, result)

	// Four of a kind check works correctly with seven cards
	hand = GenerateHand("2S 2C 2D 2H AS KC 7C")
	match, result = CheckFourOfAKind(hand)
	assert.True(t, match)
	assert.Equal(t, Tiebreakers{0}, result)

	// Four of a kind check works correctly with seven cards and extra trips
	hand = GenerateHand("5S 5C 5D 5H AS AC AH")
	match, result = CheckFourOfAKind(hand)
	assert.True(t, match)
	assert.Equal(t, Tiebreakers{3}, result)

	// No four of a kind when there is none
	hand = GenerateHand("JS JC JD 6C 6S 6D KC")
	match, result = CheckFourOfAKind(hand)
	assert.False(t, match)
	assert.Nil(t, result)
}

func TestCheckStraightFlush(t *testing.T) {
	// Straight flush check works correctly
	hand := GenerateHand("4S 5S 6S 7S 8S")
	match, result := CheckStraightFlush(hand)
	assert.True(t, match)
	assert.Equal(t, Tiebreakers{6}, result)

	// Straight flush check works correctly even with seven cards
	hand = GenerateHand("7S 8S 9S TS JS AH TC")
	match, result = CheckStraightFlush(hand)
	assert.True(t, match)
	assert.Equal(t, Tiebreakers{9}, result)

	// No straight flush if there is no flush
	hand = GenerateHand("7S 8S 9H TS JS AH TC")
	match, result = CheckStraightFlush(hand)
	assert.False(t, match)
	assert.Nil(t, result)

	// No straight flush if there is no straight
	hand = GenerateHand("6S 8S 9S TS JS AH TC")
	match, result = CheckStraightFlush(hand)
	assert.False(t, match)
	assert.Nil(t, result)
}

func TestEvaluator(t *testing.T) {
	board := "9S TS JS AH AC"
	flushAJ := "6S 8S " + board
	flushAK := "KS 8S " + board
	boatAT := "AS TH " + board
	straightA := "QS KH " + board
	handsToEvaluate := []HandForEvaluation{
		{
			Hand:        GenerateHand(flushAJ),
			PlayerIndex: 0,
		},
		{
			Hand:        GenerateHand(flushAK),
			PlayerIndex: 1,
		},
		{
			Hand:        GenerateHand(boatAT),
			PlayerIndex: 2,
		},
		{
			Hand:        GenerateHand(straightA),
			PlayerIndex: 3,
		},
	}
	evaluatedHands := EvaluateHands(handsToEvaluate)
	expectedResult := []HandForEvaluation{
		{
			Hand:         GenerateHand(boatAT),
			PlayerIndex:  2,
			HandStrength: HandStrength{7, 12, 8},
			HandName:     "Full House",
		},
		{
			Hand:         GenerateHand(flushAK),
			PlayerIndex:  1,
			HandStrength: HandStrength{6, 11, 9, 8, 7, 6},
			HandName:     "Flush",
		},
		{
			Hand:         GenerateHand(flushAJ),
			PlayerIndex:  0,
			HandStrength: HandStrength{6, 9, 8, 7, 6, 4},
			HandName:     "Flush",
		},
		{
			Hand:         GenerateHand(straightA),
			PlayerIndex:  3,
			HandStrength: HandStrength{5, 12},
			HandName:     "Straight",
		},
	}
	assert.Equal(t, expectedResult, evaluatedHands)
}
