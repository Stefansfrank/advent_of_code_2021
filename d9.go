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

}

// helper structure for points
type pnt struct {
	x int
	y int
}

// accessing the maps with index overflow safety
// basically defining the map beyond limits as 9s
func m(x,y int) int {
	if x >= xmax || y >= ymax || x < 0 || y < 0 {
		return 9
	}
	return mp[pnt{x,y}]
}

func mset(x,y,val int) {
	if x >= xmax || y >= ymax || x < 0 || y < 0 {
		return
	}
	mp[pnt{x,y}] = val
}

func mm(x,y int) bool {
	if x >= xmax || y >= ymax || x < 0 || y < 0 {
		return false
	}
	return minMp[pnt{x,y}]
}

func mmset(x,y int, val bool) {
	if x >= xmax || y >= ymax || x < 0 || y < 0 {
		return
	}
	minMp[pnt{x,y}] = val
}


// --------- P1 Main Logic -------------------


// detect the low points of the system
// NOTE: I have thought about what to do for points that are level, but:
// It turns out the puzzle input has no level points (lower than 9) next to each other 
// thus we can just omit dealing with level points. The way the code is written, 
func detMin () (mins []pnt) {

	// map indicating local mins
	// it starts out with everything being min in order to deal with margins
	minMp  = make(map[pnt]bool)
	mins   = []pnt{}

	// move sideways and down through all points comparing to the next point right / down
	// taking out possible solutions for all points compared	
	// NOTE: the getters are able to deal with index overflows thus this works ...
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

// Recursive app that looks whether the four surrounding spots are lower than 9
// and not yet tapped as part of a basin. 
// NOTE: This code does not properly deal with the special case of a basin with a sub-basin in a side wall
// if such a sub-basin would be present, two things would happen:
// a) since there are two local minimums, the basin would be counted twice
// b) since the second local minimum is already false, the basin counting would not traverse it
// a) & b) combined will lead to wonky sizes being returned for the size of the basin with sub-basin
// either one of the basins is counted with size -1 and the second one with size 1 or if the
// second minimum is a bottleneck for the traversion, both basins have bigger than 1 size but less than full size
// GOOD NEWS: no sub-basins in the input ...

func basinSize(p pnt) (size int) {

	size = 1 				// counting the point itself
	mmset(p.x,p.y, false)	// not needed for the starting local minimum but for recursion

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

// Globals (since I do not like endless signatures unless encapsilation is important)
var xmax, ymax int     // the dimensions of the map
var mp    map[pnt]int  // the map of numbers
var minMp map[pnt]bool // a map used to mark local minimums and basins (false means minimum in part 1 and depression below 9 in part 2)

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

	// Part 2
	sm := []int{}
	for _, p := range mins {
		sm = append(sm, basinSize(p))
	}
	sort.Ints(sm)
	fmt.Printf("Basin (P2): %v\n", sm[len(sm)-1] * sm[len(sm)-2] * sm[len(sm)-3])

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}