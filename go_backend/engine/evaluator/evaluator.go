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
// Returns True if the hand contains a straight, hand could be better than a straight and check_straight would still
// return true.
func CheckStraight(hand []engine.Card) (bool, Tiebreakers) {
	handRanks := getCardRankIdSlice(hand)
	descendingSort(handRanks)
	// If the highest hand rank is 12 (i.e. an ace) add a -1 to allow for
	// low end straights
	if handRanks[0] == 12 {
		handRanks = append(handRanks, -1)
	}
	for i := 0; i < len(handRanks) - 4; i++ {
		startsStraight := true
		for j := i; j < 4; j ++ {
			if handRanks[j] - handRanks[j + 1] != 1 {
				startsStraight = false
				break
			}
		}
		if startsStraight {
			return true, Tiebreakers{handRanks[i]}
		}
	}
	return false, nil
}

// Returns True if the hand contains a three of a kind, hand could be better than a three of a kind and
// check_three_of_a_kind would still return true.
func CheckThreeOfAKind(hand []engine.Card) (bool, Tiebreakers) {
	trips := getRankFrequencies(hand, 3)
	if trips == nil {
		return false, nil
	}
	trip := trips[0]
	handRanks := getCardRankIdSlice(hand)
	handRanks = filter(handRanks, func(i int) bool {
		return i != trip
	})
	descendingSort(handRanks)
	return true, append([]int{trip}, handRanks[:2]...)
}

// Returns True if the hand contains two pairs, hand could be better than two pair and check_two_pair would still
// return true.
func CheckTwoPair(hand []engine.Card) (bool, Tiebreakers) {
	pairs := getRankFrequencies(hand, 2)
	if !(len(pairs) >= 2) {
		return false, nil
	}
	pair1 := pairs[0]
	pair2 := pairs[1]
	handRanks := getCardRankIdSlice(hand)
	handRanks = filter(handRanks, func(i int) bool {
		return i != pair1 && i != pair2
	})
	descendingSort(handRanks)
	return true, []int{pair1, pair2, handRanks[0]}

}

// Returns True if the hand contains at only one pair, hand could be better than one pair and check_pair would still
// return true.
func CheckPair(hand []engine.Card) (bool, Tiebreakers) {
	pairs := getRankFrequencies(hand, 2)
	if pairs == nil {
		return false, nil
	}
	pair := pairs[0]
	handRanks := getCardRankIdSlice(hand)
	handRanks = filter(handRanks, func(i int) bool {
		return i != pair
	})
	descendingSort(handRanks)
	return true, append([]int{pair}, handRanks[:3]...)
}

// Always returns True, because hand is always at least high card good. Returns ordered cards for tie-breaking
func CheckHighCard(hand []engine.Card) (bool, Tiebreakers) {
	handRanks := getCardRankIdSlice(hand)
	descendingSort(handRanks)
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

func getRankFrequencies(hand []engine.Card, k int) []int {
	ranks := getCardRankIdSlice(hand)
	sort.Sort(sort.Reverse(sort.IntSlice(ranks)))
	var matches []int
	count := 0
	for i, current := range ranks {
		count++
		if i == len(ranks) - 1 || ranks[i + 1] != current {
			if count == k {
				matches = append(matches, current)
			}
			count = 0
		}
	}
	return matches
}

func descendingSort(vs []int) {
	sort.Slice(vs, func(i, j int) bool {
		// Use a "greater than" function instead of "less than" to sort in descending order
		return vs[i] > vs[j]
	})
}

func filter(vs []int, f func(int) bool) []int {
	vsf := make([]int, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}


