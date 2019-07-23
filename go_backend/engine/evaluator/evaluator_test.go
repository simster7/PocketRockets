package evaluator

import (
	"github.com/simster7/PocketRockets/go_backend/engine"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckHighCard(t *testing.T) {
	hand := engine.GenerateHand("8S 7H TH KD 4C")
	match, tiebreakers := CheckHighCard(hand)
	assert.True(t, match)
	assert.Equal(t, Tiebreakers{11, 8, 6, 5, 2}, tiebreakers)

	hand = engine.GenerateHand("8S 7H TH KD 4C 3C 2C")
	match, tiebreakers = CheckHighCard(hand)
	assert.True(t, match)
	assert.Equal(t, Tiebreakers{11, 8, 6, 5, 2}, tiebreakers)
}
