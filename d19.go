package main

import (
	"fmt"
	"strconv"
	"os"
	"bufio"
	"time"
	//"regexp"
	"strings"
)

// no error handling ...
func readFile (name string) (lines []string) {	
	lines = []string{}
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

// quick int abs function
func abs(x int) int {
	if (x < 0) {
		return -x
	}
	return x
}

// input parser returns the scans for all scanners
func parseFile (lines []string) (scans [][]vec) {

	for _, line := range lines {
		if len(line) > 0 {
			if line[0:3] == "---" {
				scans = append(scans, []vec{})
			} else {
				cc := strings.Split(line, ",")
				scans[len(scans)-1] = append(scans[len(scans)-1], vec{atoi(cc[0]), atoi(cc[1]), atoi(cc[2])})
			}
		}
	}
	return
} 

// basic 3D vector functions --------------------------------
type vec []int

// since we only look at a very limited set of rotations, I express the rotation as:
// - a permutation of the axes saved in 'mp'
// - a permutation of sign inversions saved in 'inv'
type rot struct {
	mp  vec
	inv vec
}

// rotates the vector sing the above representation for rotations
func (v vec) rot(r rot) (nv vec) {
	nv = vec{0,0,0}
	for i := 0; i < 3; i++ {
		nv[i] = r.inv[i] * v[r.mp[i]]
	}
	return
} 

// moves the vector 
func (v vec) trs(t vec) (nv vec) {
	return vec{v[0] + t[0], v[1] + t[1], v[2] + t[2]}
} 

// computes the movement vector needed to reach target from v
func (v vec) to(target vec) vec {
	return vec{target[0]-v[0], target[1]-v[1], target[2]-v[2]}
}

// tests for equality
func (v vec) equ(v2 vec) bool {
	return v2[0] == v[0] && v2[1] == v[1] && v2[2] == v[2]
}

// Manhattan distance
func (v vec) mdist() int {
	return abs(v[0])+abs(v[1])+abs(v[2])
}

// creates a list of allowed rotations for each scanner (permutations of axis plus negations)
// I am pretty certain, some of these are redundant and I could cut this list of 48 in half
func initRots() (rots []rot) {

	rots   = []rot{}
	inv:= []vec{ vec{ 1, 1, 1},
				 vec{ 1, 1,-1},
				 vec{ 1,-1, 1},
				 vec{-1, 1, 1},
				 vec{-1,-1, 1},
				 vec{-1, 1,-1},
				 vec{ 1,-1,-1},
				 vec{-1,-1,-1}}

	mp := []vec{ vec{ 0, 1, 2},
				 vec{ 0, 2, 1},
				 vec{ 1, 0, 2},
				 vec{ 2, 0, 1},
				 vec{ 1, 2, 0},
				 vec{ 2, 1, 0}}

	for _, i := range inv {
		for _,m := range mp {
			rots = append(rots, rot{inv:i, mp:m})
		}	
	}

	return
} 

// -------------- Core Logic --------------------------------------

// Takes the scans of two beacons and detects whether 12 or more ar matching in a valid transformation. Returns:
// - the number of matches (if greater or equal than 12)
// - a vector that contains all points of scan2 that do NOT match - transformed into the scan1 coordinate system
// - the position of the scanner in that match
func match(scan1, scan2 []vec, rots []rot) (int, []vec, vec) {

	// try all rotations
	for _, rt := range rots {

		// go through each point pair with a point from either scan and assume they are identical
		for _, p1 := range scan1 {
			for _, p2 := range scan2 {
				tr := p2.rot(rt).to(p1) // transposition to make these two points identical

				// now count how many of the points are matching now (at least one always does)
				cnt  := 0
				newP := make([]vec, 0, len(scan2) - 12)
				for _, pp2 := range scan2 {
					hit := false
					ppr2 := pp2.rot(rt).trs(tr)
					for _, pp1 := range scan1 {
						if pp1.equ(ppr2) {
							cnt += 1
							hit = true
						}						
					}
					if !hit {
						newP = append(newP, ppr2)
					}
				}

				if cnt > 11 {
					return cnt, newP, tr
				}
			}
		}
	}
	return 0, []vec{}, vec{}
}

// takes all scans and returns the total list of beacons in coordinate system of the first scanner
// it also returns a list with the coordinates of all scanners relative to scanner 0
func buildMap(scan [][]vec, rots []rot) (bcns []vec, scnrs []vec) {

 	scnrs = []vec{vec{0,0,0}}
 	bcns  = scan[0]
 	succ := map[int]bool{0:true} // keeps track which scanners are successful matched / merged into the master list of beacons

	for len(succ) < len(scan) {
		for i, sscn := range scan {

			if succ[i] {
				continue
			}

			cnt, newP, pos := match(bcns, sscn, rots)
			if cnt > 0 {
				fmt.Println("Merged", i)
				bcns = append(bcns, newP...)
				succ[i]  = true
	 			scnrs = append(scnrs, pos)
			}
		}
	}
	return
}

// cross checks all scanners with each other for maximal distance
func maxDist(scnrs []vec) (maxDst int) {
	for _, sc1 := range scnrs {
		for _, sc2 := range scnrs {
			if sc1.to(sc2).mdist() > maxDst {
				maxDst = sc1.to(sc2).mdist()
			}
		}
	}
	return
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

	start  := time.Now()

	input  := readFile("d19." + dataset + ".txt")
	
	scans  := parseFile(input)
	rots   := initRots()

	bcns, scnrs := buildMap(scans, rots)
	fmt.Println("The number of beacons is:", len(bcns))
	fmt.Println("The max distance between scanners is:", maxDist(scnrs))

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}