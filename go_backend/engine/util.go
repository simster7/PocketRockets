package engine

import "strings"

func GenerateHand(handString string) []Card {
	cardString := strings.Split(handString, " ")
	var cards []Card
	for _, card := range cardString {
		cards = append(cards, NewCard(ToCardId(string(card[0]), string(card[1]))))
	}
	return cards
}
