package engine

import "github.com/simster7/PocketRockets/go_backend/api/v1"

type Seat struct {
	Index    int
	Occupied bool
	Player   *Player
}

func (s *Seat) GetMessage() *v1.Seat {
	return &v1.Seat{
		Index:    int32(s.Index),
		Occupied: s.Occupied,
		Player:   s.Player.GetMessage(),
	}
}
