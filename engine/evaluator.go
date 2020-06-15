package engine

import (
	"sort"
)

type Tiebreakers []int
type HandStrength []int

type HandForEvaluation struct {
	Hand         []Card
	HandStrength HandStrength
	PlayerIndex  int
	HandName     string
}

func EvaluateHands(handsForEvaluation []HandForEvaluation) []HandForEvaluation {
	var evaluatedHands []HandForEvaluation
	for _, handForEval := range handsForEvaluation {
		handStrength := getHandStrength(handForEval.Hand)
		evaluatedHands = append(evaluatedHands, HandForEvaluation{
			Hand:         handForEval.Hand,
			HandStrength: handStrength,
			PlayerIndex:  handForEval.PlayerIndex,
			HandName:     getHandName(handStrength[0]),
		})
	}
	sort.Slice(evaluatedHands, func(i, j int) bool {
		// sort.Slice uses this "less" function to indicate whether i is less than j.
		// In this case, if hand i is worse than j then CompareStrengths returns 1 and -1 vice-versa.
		// However, this would sort the slice in ascending order, and we want it sorted in descending order.
		// To do that we simply negate the boolean result by checking if the result is -1 instead of 1.
		return CompareStrengths(evaluatedHands[i].HandStrength, evaluatedHands[j].HandStrength) < 0
	})
	return evaluatedHands
}

var handChecks = []func([]Card) (bool, Tiebreakers){
	CheckStraightFlush,
	CheckFourOfAKind,
	CheckFullHouse,
	CheckFlush,
	CheckStraight,
	CheckThreeOfAKind,
	CheckTwoPair,
	CheckPair,
	CheckHighCard,
}

func getHandStrength(hand []Card) HandStrength {
	possibleHands := len(handChecks)
	for i, check := range handChecks {
		match, result := check(hand)
		if match {
			handScore := possibleHands - i
			return append([]int{handScore}, result...)
		}
	}
	panic("Bug in getHandStrength")
}

func getHandName(handScore int) string {
	handNames := map[int]string{
		9: "Straight Flush",
		8: "Four Of A Kind",
		7: "Full House",
		6: "Flush",
		5: "Straight",
		4: "Three Of A Kind",
		3: "Two Pair",
		2: "Pair",
		1: "High Card",
	}
	return handNames[handScore]
}

// Returns True if the hand contains a straight flush
func CheckStraightFlush(hand []Card) (bool, Tiebreakers) {
	if len(hand) < 5 {
		return false, nil
	}
	ranks := getCardSuitIdSlice(hand)
	flush := getFrequencies(ranks, 5, 7)
	if len(flush) != 1 {
		return false, nil
	}
	suit := flush[0]
	flushMembers := filterCard(hand, func(card Card) bool {
		return card.GetSuitId() == suit
	})
	return CheckStraight(flushMembers)

}

// Returns True if the hand contains a four of a kind, hand could be better than a four of a kind and
// check_four_of_a_kind would still return true.
func CheckFourOfAKind(hand []Card) (bool, Tiebreakers) {
	if len(hand) < 4 {
		return false, nil
	}
	handRanks := getCardRankIdSlice(hand)
	quads := getFrequencies(handRanks, 4, 4)
	if len(quads) != 1 {
		return false, nil
	}
	return true, Tiebreakers{quads[0]}
}

// Returns True if the hand contains a full house, hand could be better than a full house and check_full_house would
// still return true.
func CheckFullHouse(hand []Card) (bool, Tiebreakers) {
	if len(hand) < 5 {
		return false, nil
	}
	tripMatch, tripTiebreakers := CheckThreeOfAKind(hand)
	pairMatch, pairTiebreakers := CheckPair(hand)
	if tripMatch && pairMatch {
		return true, Tiebreakers{tripTiebreakers[0], pairTiebreakers[0]}
	}
	if tripMatch {
		newHand := filterCard(hand, func(card Card) bool {
			return card.GetRankId() != tripTiebreakers[0]
		})
		tripMatch2, tripTiebreakers2 := CheckThreeOfAKind(newHand)
		if tripMatch2 {
			largerTrip, smallerTrip := maxMin(tripTiebreakers[0], tripTiebreakers2[0])
			return true, Tiebreakers{largerTrip, smallerTrip}
		}
	}
	return false, nil
}

// Returns True if the hand contains a flush, hand could be better than a flush and check_flush would still return
// true.
func CheckFlush(hand []Card) (bool, Tiebreakers) {
	if len(hand) < 5 {
		return false, nil
	}
	ranks := getCardSuitIdSlice(hand)
	flush := getFrequencies(ranks, 5, 7)
	if len(flush) != 1 {
		return false, nil
	}
	suit := flush[0]
	flushMembers := filterCard(hand, func(card Card) bool {
		return card.GetSuitId() == suit
	})
	rankedFlushMembers := getCardRankIdSlice(flushMembers)
	descendingSort(rankedFlushMembers)
	return true, rankedFlushMembers[:5]
}

// Returns True if the hand contains a straight, hand could be better than a straight and check_straight would still
// return true.
func CheckStraight(hand []Card) (bool, Tiebreakers) {
	if len(hand) < 5 {
		return false, nil
	}
	handRanks := getCardRankIdSlice(hand)
	descendingSort(handRanks)
	// If the highest hand rank is 12 (i.e. an ace) add a -1 to allow for
	// low end straights
	if handRanks[0] == 12 {
		handRanks = append(handRanks, -1)
	}
	for i := 0; i < len(handRanks)-4; i++ {
		startsStraight := true
		for j := i; j < i + 4; j++ {
			if handRanks[j]-handRanks[j+1] != 1 {
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
func CheckThreeOfAKind(hand []Card) (bool, Tiebreakers) {
	if len(hand) < 3 {
		return false, nil
	}
	ranks := getCardRankIdSlice(hand)
	trips := getFrequencies(ranks, 3, 3)
	if trips == nil {
		return false, nil
	}
	trip := trips[0]
	handRanks := getCardRankIdSlice(hand)
	handRanks = filterInt(handRanks, func(i int) bool {
		return i != trip
	})
	descendingSort(handRanks)
	return true, append([]int{trip}, handRanks[:min(2, len(handRanks))]...)
}

// Returns True if the hand contains two pairs, hand could be better than two pair and check_two_pair would still
// return true.
func CheckTwoPair(hand []Card) (bool, Tiebreakers) {
	if len(hand) < 4 {
		return false, nil
	}
	ranks := getCardRankIdSlice(hand)
	pairs := getFrequencies(ranks, 2, 2)
	if !(len(pairs) >= 2) {
		return false, nil
	}
	pair1 := pairs[0]
	pair2 := pairs[1]
	handRanks := getCardRankIdSlice(hand)
	handRanks = filterInt(handRanks, func(i int) bool {
		return i != pair1 && i != pair2
	})
	descendingSort(handRanks)
	return true, []int{pair1, pair2, handRanks[0]}

}

// Returns True if the hand contains at only one pair, hand could be better than one pair and check_pair would still
// return true.
func CheckPair(hand []Card) (bool, Tiebreakers) {
	if len(hand) < 2 {
		return false, nil
	}
	ranks := getCardRankIdSlice(hand)
	pairs := getFrequencies(ranks, 2, 2)
	if pairs == nil {
		return false, nil
	}
	pair := pairs[0]
	handRanks := getCardRankIdSlice(hand)
	handRanks = filterInt(handRanks, func(i int) bool {
		return i != pair
	})
	descendingSort(handRanks)
	return true, append([]int{pair}, handRanks[:min(3, len(handRanks))]...)
}

// Always returns True, because hand is always at least high card good. Returns ordered cards for tie-breaking
func CheckHighCard(hand []Card) (bool, Tiebreakers) {
	handRanks := getCardRankIdSlice(hand)
	descendingSort(handRanks)
	return true, handRanks[:min(5, len(handRanks))]
}

func CompareStrengths(a, b []int) int {
	for i := 0; i < len(a); i++ {
		if a[i] > b[i] {
			return -1
		} else if b[i] > a[i] {
			return 1
		}
	}
	return 0
}

func getCardRankIdSlice(hand []Card) []int {
	var handRanks []int
	for _, rank := range hand {
		handRanks = append(handRanks, rank.GetRankId())
	}
	return handRanks
}

func getCardSuitIdSlice(hand []Card) []int {
	var handSuits []int
	for _, rank := range hand {
		handSuits = append(handSuits, rank.GetSuitId())
	}
	return handSuits
}

func getFrequencies(list []int, k, j int) []int {
	sort.Sort(sort.Reverse(sort.IntSlice(list)))
	var matches []int
	count := 0
	for i, current := range list {
		count++
		if i == len(list)-1 || list[i+1] != current {
			if count >= k && count <= j {
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

func filterInt(vs []int, f func(int) bool) []int {
	vsf := make([]int, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func filterCard(vs []Card, f func(card Card) bool) []Card {
	vsf := make([]Card, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}
