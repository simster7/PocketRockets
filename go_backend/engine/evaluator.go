package engine

import "sort"

type CheckResult struct {
	Match bool
	Tiebreakers []int
}


// Returns True if the hand contains at only one pair, hand could be better than one pair and check_pair would still
// return true.
// Returns remaining cards for tie-breaking, it is structured as ([pair rank id], *[sorted kicker rank ids])
//func CheckPair(hand []Card) CheckResult {
//	return nil
//}

// Always returns True, because hand is always at least high card good. Returns ordered cards for tie-breaking
func CheckHighCard(hand []Card) CheckResult {
	var handRanks []int
	for _, rank := range hand {
		handRanks = append(handRanks, rank.GetRankId())
	}
	sort.Slice(handRanks, func(i, j int) bool {
		// Use a "greater than" function instead of "less than" to sort in descending order
		return handRanks[i] > handRanks[j]
	})
	return CheckResult{true, handRanks[:5]}
}
