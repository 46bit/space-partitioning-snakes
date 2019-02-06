package main

import (
	"encoding/json"
	"fmt"
	"github.com/46bit/circle-collision-detection/world"
	"log"
	"math"
	"math/rand"
	"time"
)

func main() {
	leftX := rand.Float64()
	topY := rand.Float64()
	bounds := world.Bounds{
		LeftX:   leftX,
		RightX:  leftX + math.Abs(rand.Float64()),
		TopY:    topY,
		BottomY: topY + math.Abs(rand.Float64()),
	}

	numberOfCircles := 100000

	start := time.Now()
	circles := make([]world.Circle, numberOfCircles)
	for i := 0; i < numberOfCircles; i++ {
		circles[i] = randomCircleWithinBounds(bounds)
	}
	elapsed := time.Now().Sub(start)
	log.Println(elapsed)

	// start = time.Now()
	// numberOfIntersections := 0
	// for i := 0; i < numberOfCircles; i++ {
	// 	for j := 0; j < numberOfCircles; j++ {
	// 		if i != j {
	// 			if circles[i].Intersects(circles[j]) {
	// 				//fmt.Printf("intersection: %d x %d\n", circles[i].ID, circles[j].ID)
	// 				numberOfIntersections += 1
	// 			}
	// 		}
	// 	}
	// }
	// log.Printf("%d of %d circles intersect\n", numberOfIntersections, numberOfCircles)
	// elapsed = time.Now().Sub(start)
	// log.Println(elapsed)

	start = time.Now()
	quadtree := world.NewQuadtree(bounds, circles)
	elapsed = time.Now().Sub(start)
	log.Println(elapsed)
	log.Printf("Circles in quadtree: %d", traverseCountCircles(quadtree))

	start = time.Now()
	n := 0
	stack := []world.Quadtree{quadtree}
	for len(stack) > 0 {
		var head world.Quadtree
		head, stack = stack[0], stack[1:]
		if head.Subtrees != nil {
			for i := range *head.Subtrees {
				//log.Println(head.Subtrees[i].Bounds)
				stack = append(stack, head.Subtrees[i])
			}
		}
		n += 1
	}
	log.Println(n)
	elapsed = time.Now().Sub(start)
	log.Println(elapsed)

	start = time.Now()
	log.Println(traverse(quadtree))
	elapsed = time.Now().Sub(start)
	log.Println(elapsed)

	start = time.Now()
	log.Println(computeComputations(quadtree, []world.Circle{}))
	elapsed = time.Now().Sub(start)
	log.Println(elapsed)

	data, err := json.MarshalIndent(quadtree, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(data))
}

func traverse(q world.Quadtree) int {
	n := 1
	if q.Subtrees != nil {
		for i := range *q.Subtrees {
			n += traverse(q.Subtrees[i])
		}
	}
	return n
}

func traverseCountCircles(q world.Quadtree) int {
	n := len(q.Circles)
	if q.Subtrees != nil {
		for i := range *q.Subtrees {
			n += traverseCountCircles(q.Subtrees[i])
		}
	}
	return n
}

func computeComputations(q world.Quadtree, parentCircles []world.Circle) int {
	n := 0

	for _, parentCircle := range parentCircles {
		for _, circle := range q.Circles {
			if parentCircle.Intersects(circle) {
				//if parentCircle.Bounds().Intersects(circle.Bounds()) && parentCircle.Intersects(circle) {
				//fmt.Printf("quadtree intersection: %d x %d\n", circle.ID, parentCircle.ID)
				//fmt.Printf("quadtree intersection: %d x %d\n", parentCircle.ID, circle.ID)
				n += 2
			}
		}
	}

	for i, circle := range q.Circles {
		for j, circle2 := range q.Circles {
			//if i != j && circle.Bounds().Intersects(circle2.Bounds()) && circle.Intersects(circle2) {
			if i != j && circle.Intersects(circle2) {
				//fmt.Printf("quadtree intersection: %d x %d\n", circle.ID, circle2.ID)
				n += 1
			}
		}
	}

	if q.Subtrees != nil {
		circles := append(parentCircles, q.Circles...)
		for i := range *q.Subtrees {
			n += computeComputations(q.Subtrees[i], circles)
		}
	}
	return n
}

func randomCircleWithinBounds(bounds world.Bounds) world.Circle {
	for {
		circle := world.Circle{
			ID: rand.Int(),
			Centre: world.Point{
				X: rand.Float64(),
				Y: rand.Float64(),
			},
			Radius: math.Abs(rand.Float64()),
		}
		if bounds.Contains(circle) {
			return circle
		}
	}
}
