package world

import (
	"math"
)

type Circle struct {
	ID     int     `json:"id"`
	Centre Point   `json:"centre"`
	Radius float64 `json:"radius"`
}

func (c Circle) Bounds() Bounds {
	return Bounds{
		LeftX:   c.Centre.X - c.Radius,
		RightX:  c.Centre.X + c.Radius,
		TopY:    c.Centre.Y - c.Radius,
		BottomY: c.Centre.Y + c.Radius,
	}
}

func (c Circle) Intersects(c2 Circle) bool {
	xDist := math.Abs(c2.Centre.X - c.Centre.X)
	yDist := math.Abs(c2.Centre.Y - c.Centre.Y)
	combinedRadii := math.Abs(c.Radius) + math.Abs(c2.Radius)
	if xDist > combinedRadii || yDist > combinedRadii {
		return false
	}
	dist := math.Hypot(xDist, yDist)
	return dist <= combinedRadii
}
