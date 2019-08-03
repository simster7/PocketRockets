package engine

type ActionType int

const (
	Check ActionType = iota
	Call
	Fold
	Bet
	Blind
)

type Action struct {
	ActionType ActionType
	Value      int
}
