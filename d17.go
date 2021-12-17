package main

import (
	"time"
	"fmt"
)

// a bounding box (target area)
type box struct {
	xfrom int
	yfrom int
	xto   int
	yto   int 
}

// a vector
type vec struct {
	x int
	y int
}

// add two vectors
func (p vec) add(v vec) vec {
	return vec{x: p.x + v.x, y: p.y + v.y}
}

// check whether a point is in a box
func (p vec) in(b box) bool {
	return p.x >= b.xfrom && p.y >= b.yfrom && p.x <= b.xto && p.y <= b.yto
}

// the puzzle input (0 my input, 1 example given)
func input() []box {
	inp := make([]box, 2)
	inp[0] = box{xfrom: 201, xto: 230, yfrom: -99, yto: -65}
	inp[1] = box{xfrom: 20, xto: 30, yfrom: -10, yto: -5}
	return inp
}

// detect whether a trajectory hits the target 
// returns the max for y as second result
func hit(vel vec, target box) (bool, int) {

	pos := vec{x: 0, y: 0}
	max := 0

	// stop when x overshoots or y is under the target
	for pos.x <= target.xto && pos.y >= target.yfrom {

		// a hit
		if pos.in(target) {
			return true, max
		}

		// next position
		pos = pos.add(vel)
		if pos.y > max {
			max = pos.y
		}

		// adapt velocity
		if vel.x > 0 {
			vel = vel.add(vec{-1,-1}) 
		} else if vel.x < 0 {
			vel = vel.add(vec{1,-1})
		} else {
			vel = vel.add(vec{0,-1})			
		}
	}

	return false, max
}

func main() {

	start := time.Now()

	inp  := input()
	cnt  := 0
	maxY := 0

	// reasonable bounds for brute force approach:
	// x needs to be below the coordinates of the target area or it overshoots on first iteration
	// y needs to be between + and - the absolute y coord of the target area
	// if smaller it overshoots on first iteration
	// if bigger it eventually comes back to zero with a velocity that overshoots on next iteration
	for x := 1; x<250; x++ {
		for y:= -100; y<100; y++ {
			hit, max := hit(vec{x,y}, inp[0])
			if hit {
				cnt += 1
				if max > maxY {
					maxY = max
				}
			}
		}
	}

	fmt.Printf("Thera are %v solutions with a maximum height of: %v\n", cnt, maxY)

	fmt.Println("Execution time:", time.Since(start))
}
