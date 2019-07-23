package engine

import "log"

var SuitMap = map[int]string {
	0: "S",
	1: "H",
	2: "C",
	3: "D",
}

var RankMap = map[int]string {
	0: "2",
	1: "3",
	2: "4",
	3: "5",
	4: "6",
	5: "7",
	6: "8",
	7: "9",
	8: "T",
	9: "J",
	10: "Q",
	11: "K",
	12: "A",
}

type Card struct {
	CardId int
}

func NewCard(cardId int) Card {
	if cardId < 0 || cardId >= 52 {
		log.Fatal("Card id must be an integer in [0, 51]")
	}
	return Card{cardId}
}

func (c *Card) GetSuitId() int {
	return c.CardId / 13
}

func (c *Card) GetSuit() string {
	return SuitMap[c.GetSuitId()]
}

func (c *Card) GetRankId() int {
	return c.CardId % 13
}

func (c *Card) GetRank() string {
	return RankMap[c.GetRankId()]
}

func ToCardId(rank, suit string) int {
	foundRank := false
	foundSuit := false
	id := 0
	for rankId, cardRank := range RankMap {
		if rank == cardRank {
			id += rankId
			foundRank = true
		}
	}
	for suitId, cardSuit := range SuitMap {
		if suit == cardSuit {
			id += suitId * 13
			foundSuit = true
		}
	}

	if !foundRank || !foundSuit {
		log.Fatal("Card is malformed")
	}
	return id
}

