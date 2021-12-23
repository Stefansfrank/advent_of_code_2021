package main

import (
	"fmt"
	"strconv"
	"os"
	"bufio"
	"time"
	"regexp"
)

// ----- some int helper functions --------------------------------------------------------------------------

// inline capable / no error handling Atoi
func atoi (x string) int {
	y, _ := strconv.Atoi(x)
	return y
}

// int abs function
func abs(x int) int {
	if (x < 0) {
		return -x
	}
	return x
}

// int abs function
func min(x,y int) int {
	if (x < y) {
		return x
	}
	return y
}

// int abs function
func max(x,y int) int {
	if (x > y) {
		return x
	}
	return y
}

// int abs function
func minmax(x,y int) (int, int) {
	if (x > y) {
		return y, x
	}
	return x, y
}

// ----- file reading and parsing ---------------------------------------------------------------------------

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

// input parser using Regex
func parseFile (lines []string) (rules []rule) {

	re  := regexp.MustCompile(`([n|f]) x=(-?[0-9]+)..(-?[0-9]+),y=(-?[0-9]+)..(-?[0-9]+),z=(-?[0-9]+)..(-?[0-9]+)`)
	tgl := map[string]bool{"n":true, "f":false}

	rules = []rule{}
	for _, line := range lines {
		match := re.FindStringSubmatch(line)

		rl   := rule{on:tgl[match[1]],bx:box{from:make([]int,3), to:make([]int, 3)}}
		for i := 0; i < 3; i++ {
			mn, mx := minmax(atoi(match[2+2*i]), atoi(match[3+2*i]))
			rl.bx.from[i] = mn
			rl.bx.to[i]   = mx
		}
		rules = append(rules, rl)
	}
	return
}

// ----- the basic structures and helpers on them ---------------------------------------------------------------------------

// a rule as parsed (bounding box and flat indicating On/Off)
type rule struct {
	on bool
	bx box
}

// executes a rule i.e. adds it's box to the space
func (r rule) exec() {
	regSet(r.bx, r.on)
}

// a bounding box
type box struct {
	from []int
	to   []int
}

// computes volume of a bounding box
func (b box) vol() int {
	return (b.to[0] - b.from[0] + 1) * (b.to[1] - b.from[1] + 1) * (b.to[2] - b.from[2] + 1)
}

// global variable representing the initialization space
var reg [][][]bool 

// sets the globes in a given space to 'tgl'
func regSet(b box, tgl bool) {
	xmax := min(b.to[0]+1, regSize+1)
	for x := max(b.from[0], -regSize); x < xmax; x++ {
		ymax := min(b.to[1]+1, regSize+1)
		for y := max(b.from[1], -regSize); y < ymax; y++ {
			zmax := min(b.to[2]+1, regSize+1)
			for z := max(b.from[2], -regSize); z < zmax; z++ {
				reg[x+regSize][y+regSize][z+regSize] = tgl
			}
		}
	}
}

// counts the lit globes
func regCnt() (cnt int) {
	xmax := regSize+1
	for x := -regSize; x < xmax; x++ {
		ymax := regSize+1
		for y := -regSize; y < ymax; y++ {
			zmax := regSize+1
			for z := -regSize; z < zmax; z++ {
				if reg[x+regSize][y+regSize][z+regSize] {
					cnt++
				}
			}
		}
	}
	return
}

// initializes the space ( ... go ... )
func regInt() {
	sz := 2*regSize+1
	reg = make([][][]bool, sz)
	for i := 0; i < sz; i++ {
		reg[i] = make([][]bool, sz)
		for j := 0; j < sz; j++ {
			reg[i][j] = make([]bool, sz)
		}
	}
}

const regSize = 50

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
	input  := readTxtFile("d22." + dataset + ".txt")
	rules  := parseFile(input)
	regInt()
	for _, rl := range rules {
		rl.exec()
	}
	fmt.Println("Total count:", regCnt())

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}