package engine

import v1 "github.com/simster7/PocketRockets/api/v1"

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
	return orderedRounds[r.GetIndex() + 1]
}

func (r Round) IsAtOrAfter(round Round) bool {
	return r.GetIndex() >= round.GetIndex()
}

type WinCondition string

const (
	WinConditionShowdown WinCondition = "Showdown"
	WinConditionFolds    WinCondition = "Folds"
)

var AllPlayers = []int{0, 1, 2, 3, 4, 5, 6, 7, 8}

type ActionConsequence struct {
	ValidAction bool
	Message     string

	PlayerIndex int
	PlayerFold  bool
	PlayerBet   int
	IsAllIn     bool

	// Refund over-bet money that can't be matched
	RefundsMoney      bool
	RefundPlayerIndex int
	RefundAmount      int

	// Ends hand
	EndsHand      bool
	Payoffs       map[Seat]int
	PotRemainder  int
	WinCondition  WinCondition
	ShowdownHands []HandForEvaluation
}


type BetVector [9]BetVectorNode

type BetVectorNode struct {
	Amount  int
	IsAllIn bool
}

type Seat struct {
	Index    int
	Occupied bool
	Player   *Player
}

func (s *Seat) GetMessage() *v1.Seat {
	if s.Occupied {
		return &v1.Seat{
			Index:    int32(s.Index),
			Occupied: true,
			Player:   s.Player.GetMessage(),
		}
	}
	return &v1.Seat{
		Index:    int32(s.Index),
		Occupied: false,
	}
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

func (a *Action) GetMessage() *v1.Action {
	return &v1.Action{
		ActionType: string(a.ActionType),
		Value:      int32(a.Value),
	}
}
