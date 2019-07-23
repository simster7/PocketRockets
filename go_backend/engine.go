package main

import (
	"fmt"
	"github.com/simster7/PocketRockets/go_backend/engine"
	"github.com/simster7/PocketRockets/go_backend/engine/evaluator"
)

func main() {
	_, sevenHigh := evaluator.CheckHighCard([]engine.Card{{1}, {2}, {3}, {4}, {5}, {6}, {7}})
	_, nineHigh := evaluator.CheckHighCard([]engine.Card{{1}, {2}, {3}, {4}, {5}, {6}, {9}})
	fmt.Println(sevenHigh)
	fmt.Println(nineHigh)
	fmt.Println(evaluator.CompareTiebreakers(nineHigh, sevenHigh))
}
