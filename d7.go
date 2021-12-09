package main

import (
	"fmt"
	"strconv"
	"os"
	"bufio"
	"time"
	"strings"
	"sort"
	"math"
)

// no error handling ...
// reading an int csv
func readCsvTxtFileInt (name string) (nums [][]int) {	
	nums     = [][]int{}
	file, _ := os.Open(name)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {	
		spl   := strings.Split(scanner.Text(), ",")
		numLn := []int{}
		for _, s := range spl {
			numLn = append(numLn, atoi(s))
		}
		nums = append(nums, numLn)
	}
	return
}

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

// cost calculation P1
func cost1 (crabs []int) (cst, lvl int) {

	// level calculation (:= median)
	sort.Ints(crabs)
	lvl = crabs[len(crabs)/2]

	// cost calculation
	cst = 0
	for _, cr := range crabs {
		cst += abs(lvl-cr)
	}
	return cst, lvl
}

// cost calculation P2 (with integrated level computation)
func cost2 (crabs []int) (cst, lvl int) {

	// level(s) calculation (:= mean)
	sum := 0
	for _, cr := range crabs {
		sum += cr
	}

	// two possible solutions since I can't figure out whether
	// I should round up/down
	lvl1 := int(math.Floor(float64(sum)/float64(len(crabs))))
	lvl2 := int(math.Ceil(float64(sum)/float64(len(crabs))))

	// calculate cost for both possible solutions
	cst1 := 0
	cst2 := 0
	for _, cr := range crabs {
		n1   := abs(lvl1-cr)
		cst1 += n1 * (n1 + 1) / 2
		n2   := abs(lvl2-cr)
		cst2 += n2 * (n2 + 1) / 2
	}

	// return the lower cost solution
	if cst2 < cst1 {
		return cst2, lvl2
	}
	return cst1, lvl1
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

	start   := time.Now()
	crabs  := readCsvTxtFileInt("d7." + dataset + ".txt")

	cst, lvl := cost1(crabs[0])
	fmt.Printf("Cost to move crabs to level %v: %v (P1 rules)\n", lvl, cst)

	cst, lvl  = cost2(crabs[0])
	fmt.Printf("Cost to move crabs to level %v: %v (P2 rules)\n", lvl, cst)

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}