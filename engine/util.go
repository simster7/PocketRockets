package engine

import (
	"strings"
)

func GenerateHand(handString string) []Card {
	cardString := strings.Split(handString, " ")
	var cards []Card
	for _, card := range cardString {
		cards = append(cards, NewCard(ToCardId(string(card[0]), string(card[1]))))
	}
	return cards
}

// An ProcessPotsPQItem is something we manage in a allInAmount queue.
type ProcessPotsPQItem struct {
	playerIndex int
	allInAmount int // The allInAmount of the item in the queue.
	index       int
}

type ProcessPotsPQ []*ProcessPotsPQItem

func (pq ProcessPotsPQ) Len() int { return len(pq) }

func (pq ProcessPotsPQ) Less(i, j int) bool {
	return pq[i].allInAmount < pq[j].allInAmount
}

func (pq ProcessPotsPQ) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *ProcessPotsPQ) Push(x interface{}) {
	n := len(*pq)
	item := x.(*ProcessPotsPQItem)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *ProcessPotsPQ) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func maxMin(a, b int) (int, int) {
	if a >= b {
		return a, b
	}
	return b, a
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func containsIntInProcessPotsPQSlice(contenders []*ProcessPotsPQItem, playerIndex int) bool {
	for _, item := range contenders {
		if item.playerIndex == playerIndex {
			return true
		}
	}
	return false
}

func containsIntInIntSlice(slice []int, i int) bool {
	for _, item := range slice {
		if item == i {
			return true
		}
	}
	return false
}
