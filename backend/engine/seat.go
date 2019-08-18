package engine

import "github.com/simster7/PocketRockets/backend/api/v1"

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
