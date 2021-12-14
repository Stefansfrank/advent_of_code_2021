package main

import (
	"fmt"
	"os"
	"bufio"
	"time"
	"regexp"
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

// input parser using Regex - returns:
// - the polymer as a map of counters for each unique letter pair
// - the rules mapping a pair to the letter to be inserted
// - the last letter which always stays the last letter but is needed for counting
func parseFile (lines []string) (poly polymer, rules map[pair]byte, last byte) {

	re   := regexp.MustCompile(`([A-Z][A-Z]) -> ([A-Z])`)
	rules = map[pair]byte{}
	poly  = map[pair]int{}

	for i, line := range lines {

		// first line
		if i == 0 {
			for j := 0; j < len(line)-1; j++ {
				poly[pair{line[j], line[j+1]}] += 1
			}
			last = line[len(line)-1]

		// third line on
		} else if i > 1 {
			match := re.FindStringSubmatch(line)
			rules[pair{match[1][0], match[1][1]}] = match[2][0]
		}
	}
	return
}

// the core logic is based on counting unique letter pairings next to each other
// this is the base structure representing a letter pair
type pair struct {
	lft byte
	rgt byte
}

// the polymer is now just a counter for each unique pair
// note that I do not remember the sequence / order of the pairs
// that would exceed my memory after 40 iterations
type polymer map[pair]int 

// iteration by creating a new polymer and counting the new pairs
func iterate(poly polymer, rules map[pair]byte, it int) (npoly map[pair]int) {

	for i := 0; i < it; i++ {
		npoly = map[pair]int{}
		for pr, cnt := range poly {
			npoly[pair{pr.lft, rules[pr]}] += cnt
			npoly[pair{rules[pr], pr.rgt}] += cnt
		}
		poly = npoly
	}

	return
}

// print out polymer (for debugging only)
func (p polymer) dump(last byte) {
	for el, cnt := range p {
		fmt.Printf("%c%c: %v\n", el.lft, el.rgt, cnt)
	}
	fmt.Printf("Last: %c\n", last)
}

// counts unique letters
// note that only the first letter in the unqiue pair counter is counted
// as the second letter is the first letter in a pair that itself is counted as well
func (p polymer) count(last byte) int {

	// counts each letter
	ix := make(map[byte]int)
	for el, cnt := range p {
		ix[el.lft] += cnt
	}
	// don't forget the last letter
	ix[last] += 1

	// create an []int of the counters for sorting
	ocs := make([]int, len(ix))
	i   := 0
	for _,cnt := range ix {
		ocs[i] = cnt
		i += 1
	}
	sort.Ints(ocs)

	return ocs[len(ocs)-1] - ocs[0]
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
	input  := readTxtFile("d14." + dataset + ".txt")
	poly, rules, last := parseFile(input)

	it := 10
	poly = iterate(poly, rules, it)	
	fmt.Printf("Highest minus lowest counter after %v iterations: %v\n", it, poly.count(last))

	it  = 30
	poly = iterate(poly, rules, it)	
	fmt.Printf("Highest minus lowest counter after another %v iterations: %v\n", it, poly.count(last))

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}