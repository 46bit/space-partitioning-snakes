package world

type Quadtree struct {
	Bounds   Bounds       `json:"bounds"`
	Circles  []Circle     `json:"circles,omitempty"`
	Subtrees *[4]Quadtree `json:"subtrees,omitempty"`
}

func NewQuadtree(bounds Bounds, circles []Circle) Quadtree {
	q := Quadtree{
		Bounds: bounds,
	}

	if len(circles) < 2 {
		if len(circles) == 1 {
			q.Circles = circles
		}
		return q
	}

	quadrants := bounds.Quadrants()
	quadrantCircles := [4][]Circle{}
	for _, circle := range circles {
		assigned := false
		for i, quadrant := range quadrants {
			if quadrant.Contains(circle) {
				quadrantCircles[i] = append(quadrantCircles[i], circle)
				assigned = true
				break
			}
		}
		if !assigned {
			q.Circles = append(q.Circles, circle)
		}
	}

	if len(quadrantCircles[0]) > 0 || len(quadrantCircles[1]) > 0 || len(quadrantCircles[2]) > 0 || len(quadrantCircles[3]) > 0 {
		q.Subtrees = &[4]Quadtree{
			NewQuadtree(quadrants[0], quadrantCircles[0]),
			NewQuadtree(quadrants[1], quadrantCircles[1]),
			NewQuadtree(quadrants[2], quadrantCircles[2]),
			NewQuadtree(quadrants[3], quadrantCircles[3]),
		}
	}

	return q
}

func (q Quadtree) Intersects(c Circle) bool {
	if !q.Bounds.Intersects(c.Bounds()) {
		return false
	}

	for _, c2 := range q.Circles {
		if c2.ID != c.ID && c2.Intersects(c) {
			return true
		}
	}

	if q.Subtrees != nil {
		for _, s := range q.Subtrees {
			if s.Intersects(c) {
				return true
			}
		}
	}

	return false
}
