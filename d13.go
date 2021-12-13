package main

import (
	"fmt"
	"strconv"
	"os"
	"bufio"
	"time"
	"strings"
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

// input parser creates:
// - a paper which has a 0 or 1 at coordinate [y][x]
// - fold instructions: fold axis in [0], direction in [1] (1:x-direction, 0:y)
func parseFile (lines []string) (pp paper, folds [][]int) {

	coords := [][]int{}
	folds   = [][]int{}
	xmax   := 0
	ymax   := 0

	for _, line := range lines {
		p := strings.Index(line, ",")
		if p > -1 {
			co := []int{atoi(line[:p]), atoi(line[p+1:])}
			if co[0] > xmax {
				xmax = co[0]
			}
			if co[1] > ymax {
				ymax = co[1]
			}
			coords = append(coords, co)
		} else {
			if len(line) > 0 {
				folds = append(folds, []int{atoi(line[13:]), 0})
				if line[11:12] == "x" {
					folds[len(folds)-1][1] = 1
				}
			}
		}
	}

	pp = make([][]int, ymax+1)
	for i := 0; i <= ymax; i++ {
		pp[i] = make([]int, xmax+1)
	}

	for _,c := range coords {
		pp[c[1]][c[0]] = 1
	}

	return pp, folds
}

// prints out visualization of paper
func (pp paper) dump() {
	for _,row := range pp {
		for _,col := range row {
			fmt.Printf("%c", rune(-11*col+46))
		}
		fmt.Println()
	}
}

// counting the dots on the paper
func (pp paper) count() (cnt int) {
	for _,row := range pp {
		for _,col := range row {
			cnt += col
		}
	}
	return
}

// the core logic of folding
// returns a new piece of paper
func (pp paper) fold(fld []int) (npp paper) {

	xd := len(pp[0])
	yd := len(pp)

	// folding horizontally
	if fld[1] == 0 {
		npp = make([][]int, fld[0])

		// map lower part up combining with content of upper part
		// NOTE: the puzzle input always folds at half line
		// thus there is no original part not covered by the fold that would
		// need to be copied to the new paper in addition to the combined part
		for y := fld[0]+1; y < yd; y++ { 
			my := fld[0] - (y - fld[0]) // mirrored coordinate
			npp[my] = make([]int, xd)
			for x := 0; x < xd; x++ {
				npp[my][x] = pp[y][x] | pp[my][x]
			}
		}

	// folding vertically
	} else { 
		npp = make([][]int, yd)
		for y := 0; y < yd; y++ {
			npp[y] = make([]int, fld[0])

			// the right part always covers the left after fold (see above)
			for x := fld[0]+1; x < xd; x++ {
				mx := fld[0] - (x - fld[0]) // mirrored coordinate
				npp[y][mx] = pp[y][x] | pp[y][mx]
			}
		}
	}

	return npp
}

type paper [][]int

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
	input  := readTxtFile("d13." + dataset + ".txt")
	paper, folds := parseFile(input)

	paper.fold(folds[0])
	fmt.Printf("Count after first fold: %v\n\nOutput:\n", paper.count())
	for _,f := range(folds) {
		paper = paper.fold(f)
	}
	paper.dump()

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}