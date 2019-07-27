package engine

type ActionType int

const (
	check ActionType = iota
	call
	fold
	bet
)

type Action struct {
	ActionType ActionType
	Value      int
}
