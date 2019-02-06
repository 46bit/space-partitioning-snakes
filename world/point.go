package world

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Bounds struct {
	LeftX   float64 `json:"left"`
	RightX  float64 `json:"right"`
	TopY    float64 `json:"top"`
	BottomY float64 `json:"down"`
}

func (b Bounds) Width() float64 {
	return b.RightX - b.LeftX
}

func (b Bounds) Height() float64 {
	return b.BottomY - b.TopY
}

func (b Bounds) Centre() Point {
	return Point{
		X: b.LeftX + b.Width()/2,
		Y: b.TopY + b.Height()/2,
	}
}

// func (b Bounds) Contains(p Point) bool {
// 	containsX := b.LeftX <= p.X && p.X < b.RightX
// 	containsY := b.TopY <= p.Y && p.Y < b.BottomY
// 	return containsX && containsY
// }

func (b Bounds) Contains(c Circle) bool {
	circleBounds := c.Bounds()
	containsX := b.LeftX <= circleBounds.LeftX && circleBounds.RightX < b.RightX
	containsY := b.TopY <= circleBounds.TopY && circleBounds.BottomY < b.BottomY
	return containsX && containsY
}

func (b Bounds) Intersects(b2 Bounds) bool {
	xOverlapWithBOnLeft := b.LeftX < b2.LeftX && b2.LeftX < b.RightX
	xOverlapWithB2OnLeft := b2.LeftX < b.LeftX && b.LeftX < b2.RightX

	yOverlapWithBOnTop := b.TopY < b2.TopY && b2.TopY < b.BottomY
	yOverlapWithB2OnTop := b2.TopY < b.TopY && b.TopY < b2.BottomY

	return (xOverlapWithBOnLeft || xOverlapWithB2OnLeft) && (yOverlapWithBOnTop || yOverlapWithB2OnTop)
}

func (b Bounds) Quadrants() [4]Bounds {
	centre := b.Centre()
	topLeft := Bounds{
		LeftX:   b.LeftX,
		RightX:  centre.X,
		TopY:    b.TopY,
		BottomY: centre.Y,
	}
	topRight := Bounds{
		LeftX:   centre.X,
		RightX:  b.RightX,
		TopY:    b.TopY,
		BottomY: centre.Y,
	}
	bottomLeft := Bounds{
		LeftX:   b.LeftX,
		RightX:  centre.X,
		TopY:    centre.Y,
		BottomY: b.BottomY,
	}
	bottomRight := Bounds{
		LeftX:   centre.X,
		RightX:  b.RightX,
		TopY:    centre.Y,
		BottomY: b.BottomY,
	}
	return [4]Bounds{topLeft, topRight, bottomLeft, bottomRight}
}
