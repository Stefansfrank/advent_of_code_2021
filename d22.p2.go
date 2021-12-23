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

// int min function
func min(x,y int) int {
	if (x < y) {
		return x
	}
	return y
}

// int max function
func max(x,y int) int {
	if (x > y) {
		return x
	}
	return y
}

// returns both min and max for ints
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

// A rule from the input (bounding box and on/off toggle)
type rule struct {
	on bool
	bx box
}

// A 3D bounding box describing a box
type box struct {
	from []int
	to   []int
}

// A bounding box with a flag indicating whether it should be added or subtracted when lights is counted
type sbox struct {
	box box
	lit bool
}

// Computes the volume of a box (careful: both from and to are part of the box)
func (b box) vol() int {
	return (b.to[0] - b.from[0] + 1) * (b.to[1] - b.from[1] + 1) * (b.to[2] - b.from[2] + 1)
}

// ----- core logic ---------------------------------------------------------------------------

// intersect logic for two boxes (the first given as an sbox)
func (b1 sbox) intersect(b2 box) box {

	intXFrom := max(b1.box.from[0], b2.from[0])
	intXTo   := min(b1.box.to[0], b2.to[0])
	intYFrom := max(b1.box.from[1], b2.from[1])
	intYTo   := min(b1.box.to[1], b2.to[1])
	intZFrom := max(b1.box.from[2], b2.from[2])
	intZTo   := min(b1.box.to[2], b2.to[2])

	if intXFrom > intXTo || intYFrom > intYTo || intZFrom > intZTo { 
		return box{from:nil, to:nil}
	}

	return box{from:[]int{intXFrom, intYFrom, intZFrom}, to:[]int{intXTo, intYTo, intZTo}}
}

// applies rules and returns a list of boxes with compensation boxes included
func boxes(rls []rule) (sbxs []sbox) {

	sbxs     = []sbox{sbox{box: rls[0].bx, lit:true}}  // the first is lit in all datasets

	// Loop over rules
	for i := 1; i < len(rls); i++ {

		// Loop over alread existing boxes
		for _,sbx := range sbxs {

			// compute intersection of new and existing boxes
			// set the signature of the intersecting box opposite to the existing box
			nsbx := sbox{box: sbx.intersect(rls[i].bx), lit: !sbx.lit}

			// test intersection
			if nsbx.box.from != nil {

				// add the compensation box
				sbxs = append(sbxs, nsbx)
			}

		}

		// add the new box if 'on' rule
		if rls[i].on {
			sbxs = append(sbxs, sbox{box: rls[i].bx, lit:true})
		}
	}
	return
}

// count a stack of boxes including compensation boxes
func countBx(sbxs []sbox) (cnt int) {
	for _, sbx := range sbxs {
		if sbx.lit {
			cnt += sbx.box.vol()
		} else {
			cnt -= sbx.box.vol()
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

	start := time.Now()
	input := readTxtFile("d22." + dataset + ".txt")
	rules := parseFile(input)

	sbxs  := boxes(rules)
	fmt.Println("Cubes Lit:",countBx(sbxs))

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}