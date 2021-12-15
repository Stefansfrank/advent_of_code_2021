package main

import (
	"fmt"
	"os"
	"bufio"
	"time"
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

// input parser 
func parseFile (lines []string) (graph map[pnt]int, dim []int) {

	graph = map[pnt]int{}
	dim   = []int{len(lines[0]), len(lines)}
	for y, line := range lines {
		for x, char := range line {
			graph[pnt{x,y}] = int(char) - '0'
		}
	}
	return
}

// simple point structure
type pnt struct {
	x int
	y int
}

// vector addition
func (p pnt) add(v pnt) pnt {
	return pnt{x:p.x+v.x, y:p.y+v.y}
}

// tests whether the point is within a matrix of dimensions dim[0], dim[1]
func (p pnt) in(dim []int) bool {
	return !(p.x < 0 || p.x >= dim[0] || p.y < 0 || p.y >= dim[1])
}

// a structure containing path relevant info for each point on the graph
type pth struct {
	visd bool // visited for analysis?
	prev pnt  // previous point on the shortest path 
			  // (only needed if the path needs to be recreatable, not strictly necessary for AoC)
	dist int  // shortest length to get to this point
}

// the core logic determining the shortest path and returning the length
// kind of a textbook Dijkstra implementation
func detShortest (graph map[pnt]int, dim []int) (len int) {

	dirs := []pnt{pnt{0,-1}, pnt{1,0}, pnt{0,1}, pnt{-1,0}}
	dijk := map[pnt]pth{}         // the matrix of path building info needed for Dijkstra
	queu := map[pnt]bool{}        // the queue is just map with points as key (value is not used)
	inf  := dim[0]*dim[1]*9 + 1   // "infinity"
	end  := pnt{dim[0]-1, dim[1]-1}

	// initialization of path building matrix with infinite distance
	for y := 0; y < dim[0]; y++ {
		for x := 0; x < dim[0]; x++ {
			dijk[pnt{x,y}] = pth{dist: inf}
		}
	}

	// start point
	cur           := pnt{0,0}         
	dijk[cur]      = pth{visd: true, dist: 0}

	finished := false
	for !finished {

		// set valid neighbours and enqueue
		for i := 0; i < 4; i++ {
			np := cur.add(dirs[i])

			// test to be inside square and not yet visited
			if np.in(dim) && !dijk[np].visd {

				// add to the queue (since using map, no double counting)
				queu[np] = true

				// if new shortest length for neighbor
				nd := dijk[cur].dist + graph[np]
				if  nd < dijk[np].dist {
					dijk[np] = pth{dist:nd, prev:cur}
				}
			}
		}

		// visit shortest path point in the queue and remove from queue
		min := inf
		for p, _ := range queu {
			if dijk[p].dist < min {
				min = dijk[p].dist
				cur = p
			}
		}
		dijk[cur] = pth{dist: dijk[cur].dist, prev: dijk[cur].prev, visd: true}
		delete(queu, cur)

		// reached the end
		if cur == end {
			finished = true
		}
	}

	return dijk[end].dist
}

// implements the Part2 modification of the input 
func expandGraph(graph map[pnt]int, dim []int) (ngraph map[pnt]int, ndim []int) {
	ndim   = []int{5*dim[0], 5*dim[1]}
	ngraph = map[pnt]int{}

	for x := 0; x < dim[1]; x++ {
		for y := 0; y < dim[1]; y++ {
			for i := 0; i < 5; i++ {
				ngraph[pnt{x+i*dim[0],y}] = add9(graph[pnt{x,y}], i)
			}
		}		
	}

	for x := 0; x < ndim[1]; x++ {
		for y := 0; y < ndim[1]; y++ {
			for i := 1; i < 5; i++ {
				ngraph[pnt{x,y+i*dim[1]}] = add9(ngraph[pnt{x,y}], i)
			}
		}		
	}

	return
}

// print the graph (for debugging)
func dump(graph map[pnt]int, dim []int) {
	for y:=0; y<dim[1]; y++ {
		for x:=0; x<dim[0]; x++ {
			fmt.Print(graph[pnt{x,y}])
		}
		fmt.Println()
	}
}

// adds 'ad' to 'i' rolling over from 9 to 1
func add9(i, ad int) int {
	return ((i + ad -1) % 9) + 1
}


// MAIN ----
func main () {

	dataset := ""
	if len(os.Args) < 2 || os.Args[1] == "" {
		fmt.Println("No argument given - trying 'test' dataset.")
		dataset = "test"
	} else {
		dataset = os.Args[1]
	}
	red  := "\033[31m"
	bold := "\033[1m"
	off  := "\033[0m"

	start  := time.Now()
	input  := readTxtFile("d15." + dataset + ".txt")
	graph, dim := parseFile(input)
 	fmt.Printf("Parsing time: %v\n", time.Since(start))
 	start   = time.Now()

	fmt.Printf("Shortest path: " + red + bold + "%v" + off + " \n", detShortest(graph, dim))
 	fmt.Printf("Path finding time (Part 1): %v\n", time.Since(start))
 	start   = time.Now()

	graph, dim  = expandGraph(graph, dim)
 	fmt.Printf("Expansion time: %v\n", time.Since(start))
 	start   = time.Now()
	fmt.Printf("Expanded path finding time (Part 2): " + red + bold + "%v" + off + " \n", detShortest(graph, dim))

 	fmt.Printf("Path finding time: %v\n", time.Since(start))
}