package engine

import "math/rand"

type StandardShuffler struct {
}

func NewStandardShuffler() *StandardShuffler {
	return &StandardShuffler{}
}

func (ss *StandardShuffler) Shuffle() Deck {
	var deck Deck
	perm := rand.Perm(52)
	for i := 0; i < 52; i++ {
		deck[perm[i]] = Card(i)
	}
	return deck
}

type DeterministicShuffler struct {
}

func NewDeterministicShuffler() *DeterministicShuffler {
	return &DeterministicShuffler{}
}

func (ds *DeterministicShuffler) Shuffle() Deck {
	var deck Deck
	for i := 0; i < 52; i++ {
		deck[i] = Card(i)
	}
	return deck
}
