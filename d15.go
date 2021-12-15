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

// input parser greates a [][]int for the graph
func parseFile (lines []string) (graph [][]int) {

	graph = make([][]int, len(lines))
	for y, line := range lines {
		graph[y] = make([]int, len(lines[0]))
		for x, char := range line {
			graph[y][x] = int(char) - '0'
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

// tests whether the point is within [0,0,xmax,ymax]
func (p pnt) in(xmax, ymax int) bool {
	return !(p.x < 0 || p.x >= xmax || p.y < 0 || p.y >= ymax)
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
func detShortest (graph [][]int) int {

	dirs := []pnt{pnt{0,-1}, pnt{1,0}, pnt{0,1}, pnt{-1,0}}
	dijk := make([][]pth, len(graph))       // the matrix of path building info needed for Dijkstra
	queu := []pnt{} 						// the queue is just a point list
	inf  := len(graph)*len(graph[0])*9 + 1  // "infinity" as a number that can't be reached
	end  := pnt{len(graph[0])-1, len(graph)-1}

	// initialization of path building matrix with infinite distance
	for y := 0; y < len(graph[1]); y++ {
		dijk[y] = make([]pth, len(graph[0]))
		for x := 0; x < len(graph[0]); x++ {
			dijk[y][x] = pth{dist: inf}
		}
	}

	// start point
	cur           		:= pnt{0,0}         
	dijk[cur.y][cur.x]   = pth{visd: true, dist: 0}

	finished := false
	for !finished {

		// set valid neighbours and enqueue
		for i := 0; i < 4; i++ {
			np := cur.add(dirs[i])

			// test to be inside square and not yet visited
			if np.in(len(graph[0]), len(graph)) && !dijk[np.y][np.x].visd {

				// if that neighbour has never been enqueued (dist == inf) - do so
				if dijk[np.y][np.x].dist == inf {
					queu = append(queu, np)					
				}

				// is the dist to the neighbor new shortest length to it?
				nd := dijk[cur.y][cur.x].dist + graph[np.y][np.x]
				if  nd < dijk[np.y][np.x].dist {
					dijk[np.y][np.x] = pth{dist:nd, prev:cur}
				}
			}
		}

		// visit shortest path point in the queue and remove it from queue
		min := inf
		ix  := -1
		for i, p := range queu {
			if dijk[p.y][p.x].dist < min {
				min = dijk[p.y][p.x].dist
				cur = p
				ix  = i
			}
		}
		dijk[cur.y][cur.x].visd = true
		queu = append(queu[:ix], queu[ix+1:]...)

		// reached the end
		if cur == end {
			finished = true
		}
	}

	return dijk[end.y][end.x].dist
}

// implements the Part2 modification of the input 
func expandGraph(graph [][]int) (ngraph [][]int) {
	ngraph = make([][]int, 5*len(graph))

	for y := 0; y < len(graph); y++ {
		ngraph[y] = make([]int, 5*len(graph[0]))
		for x := 0; x < len(graph[0]); x++ {
			for i := 0; i < 5; i++ {
				ngraph[y][x+i*len(graph[0])] = add9(graph[y][x], i)
			}
		}		
	}

	for y := 0; y < len(graph); y++ {
		for i := 1; i < 5; i++ {
			ngraph[y+i*len(graph)] = make([]int, 5*len(graph[0]))
			for x := 0; x < 5*len(graph[0]); x++ {
				ngraph[y+i*len(graph)][x] = add9(ngraph[y][x], i)
			}
		}		
	}
	return
}

// print the graph (for debugging)
func dump(graph [][]int) {
	for y:=0; y<len(graph); y++ {
		for x:=0; x<len(graph[0]); x++ {
			fmt.Print(graph[y][x])
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
	graph  := parseFile(input)
 	fmt.Printf("Parsing: %v\n", time.Since(start))
 	start   = time.Now()

	fmt.Printf("Shortest path: " + red + bold + "%v" + off + "\n", detShortest(graph))
 	fmt.Printf("Part1: %v\n", time.Since(start))
 	start   = time.Now()

	graph   = expandGraph(graph)
 	fmt.Printf("Expand cave: %v\n", time.Since(start))
 	start   = time.Now()
	fmt.Printf("Shortest path expanded cave:" + red + bold + "%v" + off + "\n", detShortest(graph))

 	fmt.Printf("Part2: %v\n", time.Since(start))
}