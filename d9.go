package main

import (
	"fmt"
	"strconv"
	"os"
	"bufio"
	"time"
	"sort"
)

// no error handling ...
func readTxtFile (name string) (lines []string) {	
	file, _ := os.Open(name)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {		
		lines = append(lines, scanner.Text())
	}
	return
}

// inline capable / no error handling Atoi
func atoi (x string) int {
	y, _ := strconv.Atoi(x)
	return y
}

// parsing input into the global map (map[pnt]int)
func parseFile (lines []string) {
	xmax = len(lines[0])
	ymax = len(lines)
	for y, line := range lines {
		for x, _ := range line {
			mset(x,y,atoi(line[x:x+1]))
		}
	}

	// padding
	for x := -1; x <= xmax; x++ {
		mset(x,-1,9)
		mset(x,ymax,9)
		mmset(x,-1,false)
		mmset(x,ymax,false)
	}
	for y := 0; y <= ymax; y++ {
		mset(-1,y,9)
		mset(xmax,y,9)
		mmset(-1,y,false)
		mmset(xmax,y,false)
	}

}

// helper structure for points
type pnt struct {
	x int
	y int
}

// Globals (since I am lazy and do not like endless signatures unless encapsilation is important)
var xmax, ymax int     // the dimensions of the map
var mp    map[pnt]int  // the map of numbers
var minMp map[pnt]bool // a masking map used to identify local minimums and basins (Note: false means loc. minimum / basin)

// helpers to make code shorter / more readable
func m(x,y int) int {
	return mp[pnt{x,y}]
}

func mset(x,y,val int) {
	mp[pnt{x,y}] = val
}

func mm(x,y int) bool {
	return minMp[pnt{x,y}]
}

func mmset(x,y int, val bool) {
	minMp[pnt{x,y}] = val
}

// --------- P1 Main Logic -------------------


// detect the low points of the system
// NOTE: I have thought about what to do for points that are level, but:
// It turns out the puzzle input has no level minimums next to each other 
// thus we can just omit dealing with level points when searching for the minimum 
func detMin () (mins []pnt) {

	// map indicating local mins
	// it starts out with everything being min in order to deal with margins
	minMp  = make(map[pnt]bool)
	mins   = []pnt{}

	// move sideways and down through all points comparing to the next point right / down
	// taking out possible solutions for all points compared	
	for x := 0; x < xmax; x++ {
		for y :=0; y < ymax; y++ {

			// sideways move
			if m(x, y) < m(x+1, y) {
				mmset(x+1, y, true)
			} else if m(x, y) > m(x+1, y) {
				mmset(x, y, true)
			} else {
				mmset(x+1, y, true)
				mmset(x, y, true)
			}

			// downward move
			if m(x,y) < m(x,y+1) {
				mmset(x, y+1, true)
			} else if m(x,y) > m(x,y+1) {
				mmset(x, y, true)
			} else {
				mmset(x, y+1, true)
				mmset(x, y, true)				
			}
		}
	}

	// collects all points that have not been marked excluded by the loops above
	for x := 0; x < xmax; x++ {
		for y := 0; y < ymax; y++ {
			if !mm(x,y) {
				mins = append(mins, pnt{x,y})
			}
		}
	}

	return
}

// result computation for P1
func sum (vs []pnt) (sm int) {
	for _, v := range vs {
		sm += m(v.x, v.y)
	}
	return
}

// ------------- P2 Main Logic ----------------------

// Recursive app that starts at a minimum and adds neighbouring tiles unless they are 9. 
// NOTE: This code does not properly deal with the special case of a basin with a sub-basin in a side wall
// if such a sub-basin would be present, one basin would be split in two basins in the result 
// However, the puzzle input does not have any such occurance

func basinSize(p pnt) (size int) {

	size = 1 				// counting the point itself
	mmset(p.x,p.y, false)	// not needed for the starting local minimum (already false) but for recursion

	// loop through the four neighbors recursively
	if m(p.x+1,p.y) < 9 && mm(p.x+1,p.y) {
		size += basinSize(pnt{p.x+1,p.y})
	}
	if m(p.x-1,p.y) < 9 && mm(p.x-1,p.y) {
		size += basinSize(pnt{p.x-1,p.y})
	}
	if m(p.x,p.y+1) < 9 && mm(p.x,p.y+1) {
		size += basinSize(pnt{p.x,p.y+1})
	}
	if m(p.x,p.y-1) < 9 && mm(p.x,p.y-1) {
		size += basinSize(pnt{p.x,p.y-1})
	}

	return
}

// MAIN ----
func main () {

	// command line input to switch data sets
	dataset := ""
	if len(os.Args) < 2 || os.Args[1] == "" {
		fmt.Println("No argument given - trying 'test' dataset.")
		dataset = "test"
	} else {
		dataset = os.Args[1]
	}

	start  := time.Now()
	mp      = make(map[pnt]int)
	minMp   = make(map[pnt]bool)
	input  := readTxtFile("d9." + dataset + ".txt")
	parseFile(input) // global map is built

	// Part 1
	mins := detMin() // min map is built
	fmt.Printf("Risk  (P1): %v\n", sum(mins)+len(mins))
 	fmt.Printf("Execution time: %v\n", time.Since(start))
	start   = time.Now()

	// Part 2
	sm := []int{}
	for _, p := range mins {
		sm = append(sm, basinSize(p))
	}
	sort.Ints(sm)
	fmt.Printf("Basin (P2): %v\n", sm[len(sm)-1] * sm[len(sm)-2] * sm[len(sm)-3])

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}