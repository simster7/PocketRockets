package engine

import "github.com/simster7/PocketRockets/go_backend/api"

type Seat struct {
	Index    int
	Occupied bool
	Player   *Player
}

func (s *Seat) GetMessage() *api.Seat {
	return &api.Seat{
		Index:    int32(s.Index),
		Occupied: s.Occupied,
		Player:   s.Player.GetMessage(),
	}
}
