package main

import (
	//"encoding/json"
	"fmt"
	"github.com/46bit/circle-collision-detection/world"
	"log"
	"math"
	"math/rand"
	"time"
)

//var randScaleFactor = 200.0

func main() {
	rand.Seed(time.Now().UnixNano())

	// leftX := rand.Float64() * randScaleFactor
	// topY := rand.Float64() * randScaleFactor
	// bounds := world.Bounds{
	// 	LeftX:   leftX,
	// 	RightX:  leftX + math.Abs(rand.Float64()*randScaleFactor),
	// 	TopY:    topY,
	// 	BottomY: topY + math.Abs(rand.Float64()*randScaleFactor),
	// }
	bounds := world.Bounds{
		LeftX:   -100,
		RightX:  100,
		TopY:    -100,
		BottomY: 100,
	}

	start := time.Now()
	numberOfCircles := 0
	circles := make([]world.Circle, numberOfCircles)
	for i := 0; i < numberOfCircles; i++ {
		circles[i] = randomCircleWithinBounds(bounds)
	}
	elapsed := time.Now().Sub(start)
	log.Println(elapsed)

	start = time.Now()
	numberOfSnakes := 75
	snakes := make([][]world.Circle, numberOfSnakes)
	for i := 0; i < numberOfSnakes; i++ {
		snakes[i] = randomSnake(uint(rand.Int63n(60)), bounds)
		for j := 0; j < len(snakes[i]); j++ {
			circles = append(circles, snakes[i][j])
		}
	}
	elapsed = time.Now().Sub(start)
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

	start = time.Now()
	sampleCircles := [500]world.Circle{}
	for i := 0; i < 500; i++ {
		sampleCircles[i] = randomCircleWithinBounds(bounds)
		//sampleCircles[i].Radius /= 1000000
	}
	elapsed = time.Now().Sub(start)
	log.Println(elapsed)

	start = time.Now()
	sampleCirclesIntersectingCount := 0
	for i := 0; i < 500; i++ {
		if quadtree.Intersects(sampleCircles[i]) {
			sampleCirclesIntersectingCount += 1
		}
	}
	log.Println(sampleCirclesIntersectingCount)
	elapsed = time.Now().Sub(start)
	log.Println(elapsed)

	// for i := 0; i < 100000; i++ {
	// 	// if i%1 == 0 {
	// 	// log.Printf("i=%d\n", i)
	// 	// }
	// 	c := randomCircleWithinBounds(bounds)
	// 	if !quadtree.Intersects(c) {
	// 		log.Printf("No intersection: %#v", c)
	// 		//break
	// 	}
	// }

	// data, err := json.MarshalIndent(quadtree, "", "  ")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(string(data))

	fmt.Println(`<svg viewBox="0 0 200 200" xmlns="http://www.w3.org/2000/svg">`)
	for _, c := range circles {
		fmt.Printf(`  <circle cx="%f" cy="%f" r="%f" style="fill: black; stroke-width: 0.01;" />`, c.Centre.X-bounds.LeftX, c.Centre.Y-bounds.TopY, c.Radius)
		fmt.Println()
	}
	traversePrintSquares(quadtree, bounds, 1)
	fmt.Println(`</svg>`)
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

func traversePrintSquares(q world.Quadtree, bounds world.Bounds, depth int) {
	fmt.Printf(
		`  <rect x="%f" y="%f" width="%f" height="%f" style="fill: transparent; stroke: red; stroke-width: %f;" />`,
		q.Bounds.LeftX-bounds.LeftX,
		q.Bounds.TopY-bounds.TopY,
		q.Bounds.Width(),
		q.Bounds.Height(),
		1/float64(depth),
	)
	fmt.Println()
	if q.Subtrees != nil {
		for _, t := range *q.Subtrees {
			traversePrintSquares(t, bounds, depth+1)
		}
	}
}

func computeComputations(q world.Quadtree, parentCircles []world.Circle) int {
	n := 0

	for _, parentCircle := range parentCircles {
		for _, circle := range q.Circles {
			if parentCircle.Intersects(circle) {
				//fmt.Printf("quadtree intersection: %d x %d\n", circle.ID, parentCircle.ID)
				//fmt.Printf("quadtree intersection: %d x %d\n", parentCircle.ID, circle.ID)
				n += 2
			}
		}
	}

	for i, circle := range q.Circles {
		for j, circle2 := range q.Circles {
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
				X: (rand.Float64() - 0.5) * bounds.Width(),
				Y: (rand.Float64() - 0.5) * bounds.Height(),
			},
			Radius: math.Abs(rand.Float64()*math.Min(bounds.Width(), bounds.Height())) / 200,
		}
		if bounds.Contains(circle) {
			return circle
		}
	}
}

func randomSnake(length uint, bounds world.Bounds) []world.Circle {
	head := randomCircleWithinBounds(bounds)
	distance := rand.Float64() * head.Radius * 1.5
	angle := rand.Float64() * math.Pi

	segments := []world.Circle{head}
	for i := 0; i < int(length)-1; i++ {
		offset := world.Point{
			X: distance * math.Cos(angle),
			Y: distance * math.Sin(angle),
		}
		angle += 0.3 * (rand.Float64() - 0.5)

		segment := world.Circle{
			ID: rand.Int(),
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
	return segments
}
