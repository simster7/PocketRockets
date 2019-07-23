package evaluator

import (
	"github.com/simster7/PocketRockets/go_backend/engine"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetRankFrequencies(t *testing.T) {
	hand := engine.GenerateHand("8S 7H 7S 7D 3S 3C 2C 2S")
	pairs := getRankFrequencies(hand, 2)
	assert.Equal(t, []int{1, 0}, pairs)

	hand = engine.GenerateHand("8S 7H 7S 7D 3S KC 9C 2S")
	pairs = getRankFrequencies(hand, 2)
	assert.Nil(t, pairs)

	hand = engine.GenerateHand("8S 7H 7S 7D 3S 3C 2C 2S")
	trips := getRankFrequencies(hand, 3)
	assert.Equal(t, []int{5}, trips)

	hand = engine.GenerateHand("7C 7H 7S 7D 3S 3C 2C 2S")
	fours := getRankFrequencies(hand, 4)
	assert.Equal(t, []int{5}, fours)
}

func TestCheckHighCard(t *testing.T) {
	// High card evaluated correctly
	hand := engine.GenerateHand("8S 7H TH KD 4C")
	match, tiebreakers := CheckHighCard(hand)
	assert.True(t, match)
	assert.Equal(t, Tiebreakers{11, 8, 6, 5, 2}, tiebreakers)

	// High card evaluated correctly even for seven cards
	hand = engine.GenerateHand("8S 7H TH KD 4C 3C 2C")
	match, tiebreakers = CheckHighCard(hand)
	assert.True(t, match)
	assert.Equal(t, Tiebreakers{11, 8, 6, 5, 2}, tiebreakers)
}

func TestCheckPair(t *testing.T) {
	// One pair correctly evaluated
    hand := engine.GenerateHand("8S 8H TH KD 4C")
    match, result := CheckPair(hand)
    assert.True(t, match)
    assert.Equal(t, Tiebreakers{6, 11, 8, 2}, result)

    // One pair correctly evaluated with seven cards
	hand = engine.GenerateHand("8S 8H KH JD 4C AH 9C")
	match, result = CheckPair(hand)
    assert.True(t, match)
    assert.Equal(t, Tiebreakers{6, 12, 11, 9}, result)

	// No pair not evaluated, even with seven cards
	hand = engine.GenerateHand("8S 7H KH TD QC 2S 4H")
	match, result = CheckPair(hand)
    assert.False(t, match)
    assert.Nil(t, result)
}

func TestCheckTwoPair(t *testing.T) {
	// Two pair correctly evaluated
    hand := engine.GenerateHand("8S 8H TH TD 4C")
    match, result := CheckTwoPair(hand)
    assert.True(t, match)
    assert.Equal(t, Tiebreakers{8, 6, 2}, result)

    // Two pair correctly evaluated with seven cards
    hand = engine.GenerateHand("8S 8H TH TD 4C 6C 9S")
    match, result = CheckTwoPair(hand)
    assert.True(t, match)
    assert.Equal(t, Tiebreakers{8, 6, 7}, result)

    // Two pair correctly evaluated with seven cards and presence of three pair
    hand = engine.GenerateHand("8S 8H TH TD 9C 9S 6C")
    match, result = CheckTwoPair(hand)
    assert.True(t, match)
    assert.Equal(t, Tiebreakers{8, 7, 6}, result)

    // No two pair evaluated, even with seven cards
    hand = engine.GenerateHand("8S 7H KH TD QC 2S 4H")
    match, result = CheckTwoPair(hand)
    assert.False(t, match)
    assert.Nil(t, result)

    // No two pair evaluated, even with seven cards and a single pair
    hand = engine.GenerateHand("8S 8H KH TD QC 2S 4H")
    match, result = CheckTwoPair(hand)
    assert.False(t, match)
    assert.Nil(t, result)
}

func TestCheckThreeOfAKind(t *testing.T) {
	// Three of a kind correctly evaluated
    hand := engine.GenerateHand("8S 8H 8C TD 4C")
    match, result := CheckThreeOfAKind(hand)
    assert.True(t, match)
    assert.Equal(t, Tiebreakers{6, 8, 2}, result)

    // Three of a kind correctly evaluated, even with seven cards
    hand = engine.GenerateHand("8S 8H 8C TD 4C JC 2S")
    match, result = CheckThreeOfAKind(hand)
    assert.True(t, match)
    assert.Equal(t, Tiebreakers{6, 9, 8}, result)

    // No three of a kind evaluated, even with seven cards
    hand = engine.GenerateHand("8S 8H 7C TD 4C JC 2S")
    match, result = CheckThreeOfAKind(hand)
    assert.False(t, match)
    assert.Nil(t, result)

    // Three of a kind correctly evaluated, even with seven cards and extra pair (full house)
    // This tests undefined behavior, so if it fails it is safe to delete
    hand = engine.GenerateHand("8S 8H 8C TD TC JC 2S")
    match, result = CheckThreeOfAKind(hand)
    assert.True(t, match)
    assert.Equal(t, Tiebreakers{6, 9, 8}, result)
}

func TestCheckStraight(t *testing.T) {
	// Straight check works correctly
    hand := engine.GenerateHand("2S 3H 4C 5D 6C")
    match, result := CheckStraight(hand)
    assert.True(t, match)
    assert.Equal(t, Tiebreakers{4}, result)

    // Straight check works correctly, even with seven cards
    hand = engine.GenerateHand("2S 3H 4C 5D 6C KC JS")
    match, result = CheckStraight(hand)
    assert.True(t, match)
    assert.Equal(t, Tiebreakers{4}, result)

    // traight check works correctly, even with seven cards and straight is more than five cards
    hand = engine.GenerateHand("2S 3H 4C 5D 6C 7C JS")
    match, result = CheckStraight(hand)
    assert.True(t, match)
    assert.Equal(t, Tiebreakers{5}, result)

    // Straight check works correctly, even with seven cards and is ace low
    hand = engine.GenerateHand("2S 3H 4C 5D AC 9C JS")
    match, result = CheckStraight(hand)
    assert.True(t, match)
    assert.Equal(t, Tiebreakers{3}, result)

    // Straight check works correctly, even with seven cards and is ace high
    hand = engine.GenerateHand("TS QH KC 5D AC 9C JS")
    match, result = CheckStraight(hand)
    assert.True(t, match)
    assert.Equal(t, Tiebreakers{12}, result)

    // No straight when there is not straight
    hand = engine.GenerateHand("8S QH KC 6D AC 9C JS")
    match, result = CheckStraight(hand)
    assert.False(t, match)
    assert.Nil(t, result)

    // No straight wrap-around
	hand = engine.GenerateHand("8S QH KC 6D AC 2C JS")
	match, result = CheckStraight(hand)
	assert.False(t, match)
	assert.Nil(t, result)
}
