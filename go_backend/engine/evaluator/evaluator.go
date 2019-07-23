package evaluator

import (
	"github.com/simster7/PocketRockets/go_backend/engine"
	"sort"
)

type Tiebreakers []int

func CompareTiebreakers(a, b Tiebreakers) int {
	for i := 0; i < len(a); i++ {
		if a[i] > b[i] {
			return -1
		} else if b[i] > a[i] {
			return 1
		}
	}
	return 0
}



// Returns True if the hand contains at only one pair, hand could be better than one pair and check_pair would still
// return true.
// Returns remaining cards for tie-breaking, it is structured as ([pair rank id], *[sorted kicker rank ids])
func CheckPair(hand []Card) (bool, Tiebreakers) {
	handRanks := getCardRankIdSlice(hand)
	return false, nil
}

// Always returns True, because hand is always at least high card good. Returns ordered cards for tie-breaking
func CheckHighCard(hand []engine.Card) (bool, Tiebreakers) {
	handRanks := getCardRankIdSlice(hand)
	sort.Slice(handRanks, func(i, j int) bool {
		// Use a "greater than" function instead of "less than" to sort in descending order
		return handRanks[i] > handRanks[j]
	})
	return true, handRanks[:5]
}

func getCardRankIdSlice(hand []engine.Card) []int {
	var handRanks []int
	for _, rank := range hand {
		handRanks = append(handRanks, rank.GetRankId())
	}
	return handRanks
}

func getCardSuitIdSlice(hand []engine.Card) []int {
	var handSuits []int
	for _, rank := range hand {
		handSuits = append(handSuits, rank.GetSuitId())
	}
	return handSuits
}

func getRankFrequencies(hand []engine.Card) map[int][]int {
	var frequencies map[int][]int
	for _
}

