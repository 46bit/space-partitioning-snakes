package main

import (
	//"encoding/json"
	"fmt"
	"github.com/46bit/space-partitioning-snakes/world"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	bounds := world.Bounds{
		LeftX:   -1000,
		RightX:  1000,
		TopY:    -1000,
		BottomY: 1000,
	}

	start := time.Now()
	numberOfSnakes := 500
	snakes := map[int]world.Snake{}
	velocities := map[int]world.Velocity{}
	for i := 0; i < numberOfSnakes; i++ {
		snake, velocity := randomSnake(uint(rand.Int63n(160)), bounds)
		snakes[snake.ID] = snake
		velocities[snake.ID] = velocity
	}
	log.Printf("Random snakes and velocities: %s\n", time.Now().Sub(start))

	for f := 0; true; f++ {
		start = time.Now()
		circles := []world.Circle{}
		for _, snake := range snakes {
			circles = append(circles, snake.Head)
			circles = append(circles, snake.Tail...)
		}
		log.Printf("Snakes into circles: %s\n", time.Now().Sub(start))

		start = time.Now()
		quadtree := world.NewQuadtree(bounds, circles)
		log.Printf("Circle into quadtree: %s\n", time.Now().Sub(start))

		start = time.Now()
		snakes = collisions(snakes, quadtree)
		log.Printf("Snake collisions and deaths: %s\n", time.Now().Sub(start))

		if f%10 == 0 {
			start = time.Now()
			for l := 0; l < 20; l++ {
				s := `<svg viewBox="0 0 2000 2000" xmlns="http://www.w3.org/2000/svg">`
				s += traversePrintAll(quadtree, bounds, 1, l)
				s += `</svg>`
				err := ioutil.WriteFile(fmt.Sprintf("levels-%d.svg", l), []byte(s), 0644)
				if err != nil {
					fmt.Errorf(err.Error())
				}
			}
			log.Printf("Drawing SVGs: %s\n", time.Now().Sub(start))

			time.Sleep(1000 * time.Millisecond)
		}

		start = time.Now()
		for id := range snakes {
			snakes[id] = world.ApplyVelocityToSnake(velocities[id], snakes[id], false)
			velocity := velocities[id]
			velocity.Angle += 0.3 * (rand.Float64() - 0.5)
			velocities[id] = velocity
		}
		log.Printf("Apply velocities to snakes: %s\n", time.Now().Sub(start))
	}
}

func traversePrintAll(q world.Quadtree, bounds world.Bounds, depth, maxDepth int) string {
	if maxDepth < 1 {
		return ""
	}
	s := fmt.Sprintf(
		`  <rect x="%f" y="%f" width="%f" height="%f" style="fill: transparent; stroke: red; stroke-width: %f;" />`,
		q.Bounds.LeftX-bounds.LeftX,
		q.Bounds.TopY-bounds.TopY,
		q.Bounds.Width(),
		q.Bounds.Height(),
		1/float64(depth),
	)
	s += "\n"
	for _, c := range q.Circles {
		s += fmt.Sprintf(
			`  <circle cx="%f" cy="%f" r="%f" style="fill: black; stroke-width: 0.01;" />`,
			c.Centre.X-bounds.LeftX,
			c.Centre.Y-bounds.TopY,
			c.Radius,
		)
		s += "\n"
	}
	if q.Subtrees != nil {
		for _, t := range *q.Subtrees {
			s += traversePrintAll(t, bounds, depth+1, maxDepth-1)
		}
	}
	return s
}

func randomCircleWithinBounds(bounds world.Bounds) world.Circle {
	for {
		radius := 0.4
		circle := world.Circle{
			ID: rand.Int(),
			Centre: world.Point{
				X: (rand.Float64() - 0.5) * bounds.Width(),
				Y: (rand.Float64() - 0.5) * bounds.Height(),
			},
			Radius: radius,
		}
		if bounds.Contains(circle) {
			return circle
		}
	}
}

func randomSnake(length uint, bounds world.Bounds) (world.Snake, world.Velocity) {
	head := randomCircleWithinBounds(bounds)
	distance := head.Radius * 1.5
	angle := rand.Float64() * math.Pi
	headAngle := angle

	segments := []world.Circle{head}
	for i := 0; i < int(length)-1; i++ {
		offset := world.Point{
			X: distance * math.Cos(math.Pi-angle),
			Y: distance * math.Sin(-angle),
		}
		angle += 0.3 * (rand.Float64() - 0.5)

		segment := world.Circle{
			ID: head.ID,
			Centre: world.Point{
				X: segments[i].Centre.X + offset.X,
				Y: segments[i].Centre.Y + offset.Y,
			},
			Radius: head.Radius * 0.95,
		}
		if !bounds.Contains(segment) {
			return randomSnake(length, bounds)
		}
		segments = append(segments, segment)
	}
	return world.Snake{
			ID:   head.ID,
			Head: head,
			Tail: segments[1:],
		}, world.Velocity{
			Speed: distance,
			Angle: headAngle,
		}
}

func collisions(snakes map[int]world.Snake, quadtree world.Quadtree) map[int]world.Snake {
	result := make(map[int]world.Snake, len(snakes))
	var mut sync.Mutex
	var wg sync.WaitGroup
	for id, snake := range snakes {
		wg.Add(1)
		go func(id int, snake world.Snake) {
			if !quadtree.Bounds.Intersects(snake.Head.Bounds()) || quadtree.Intersects(snake.Head) {
				log.Printf("Snake %d died\n", snake.ID)
			} else {
				mut.Lock()
				result[id] = snake
				mut.Unlock()
			}
			wg.Done()
		}(id, snake)
	}
	wg.Wait()
	return result
}
