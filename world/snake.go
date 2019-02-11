package world

type Snake struct {
	ID   int      `json:"id"`
	Head Circle   `json:"head"`
	Tail []Circle `json:"segments"`
}

func (s Snake) Collided(q Quadtree) bool {
	return q.Intersects(s.Head)
}
