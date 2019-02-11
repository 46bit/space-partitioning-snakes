package world

import (
	"math"
)

type Velocity struct {
	Speed float64 `json:"speed"`
	Angle float64 `json:"angle"`
}

func ApplyVelocityToSnake(v Velocity, s Snake, grow bool) Snake {
	if !grow && len(s.Tail) > 0 {
		s.Tail = s.Tail[:len(s.Tail)-1]
	}
	s.Tail = append([]Circle{s.Head}, s.Tail...)

	offset := Point{
		X: v.Speed * math.Cos(v.Angle),
		Y: v.Speed * math.Sin(v.Angle),
	}
	s.Head.Centre = Point{
		X: s.Head.Centre.X + offset.X,
		Y: s.Head.Centre.Y + offset.Y,
	}

	return s
}
