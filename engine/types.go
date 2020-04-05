package engine

type Deck [52]Card

type FoldVector [9]bool

type Round string

const (
	RoundPreFlop Round = "PreFlop"
	RoundFlop    Round = "Flop"
	RoundTurn    Round = "Turn"
	RoundRiver   Round = "River"
	RoundHandEnd Round = "HandEnd"
)

var orderedRounds = []Round{RoundPreFlop, RoundFlop, RoundTurn, RoundRiver, RoundHandEnd}

func (r Round) GetIndex() int {
	for i, v := range orderedRounds {
		if r == v {
			return i
		}
	}
	// Unreachable
	return -1
}

func (r Round) GetNextRound() Round {
	if r == RoundHandEnd {
		return r
	}
	return orderedRounds[r.GetIndex()+1]
}

func (r Round) IsAtOrAfter(round Round) bool {
	return r.GetIndex() >= round.GetIndex()
}

var AllPlayers = []int{0, 1, 2, 3, 4, 5, 6, 7, 8}

type Players [9]*Player

type Player struct {
	Name       string
	Stack      int
	Bet        int
	Folded     bool
	IsAllIn    bool
	LastAction Action
	SittingOut bool
}

type Seats [9]*Seat

type Seat struct {
	Name       string
	Stack      int
	SittingOut bool
}

type ActionType string

const (
	ActionTypeCheck ActionType = "Check"
	ActionTypeCall  ActionType = "Call"
	ActionTypeFold  ActionType = "Fold"
	ActionTypeBet   ActionType = "Bet"
	ActionTypeBlind ActionType = "Blind"
)

type Action struct {
	ActionType ActionType
	Value      int
}
