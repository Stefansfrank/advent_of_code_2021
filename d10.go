package main

import (
	"fmt"
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

// main logic for both parts
func navRepair(nav []string) (brokSum int, midScore int) {

	scores := []int{} // scores for completion of incomplete lines

	lineLoop: for _, n := range nav {
		lifo := []int{} // the lifo of open brackets by type index (0,1,2,3)

		for _,b := range n {
			if boc[b] > 0 {
				lifo = append(lifo, bix[b])
			} else {
				if lifo[len(lifo)-1] == bix[b] {
					lifo = lifo[:len(lifo)-1]
				} else {
					brokSum += bpt[b]
					continue lineLoop				
				}
			}
		}

		sc := 0
		for i := len(lifo)-1; i >= 0; i-- {
			sc = sc * 5 + lifo[i] + 1 // the indexing of bracket types happens to be the points-1
		}
		scores = append(scores, sc)
	}

	sort.Ints(scores)
	return brokSum, scores[len(scores)/2]
}

// puzzle rules init
func ini() {
	bix = map[rune]int{'(':0, ')':0, '[':1, ']':1, '{':2, '}':2, '<':3, '>':3}
	boc = map[rune]int{'(':1, ')':-1, '[':1, ']':-1, '{':1, '}':-1, '<':1, '>':-1}
	bpt = map[rune]int{'(':3, ')':3, '[':57, ']':57, '{':1197, '}':1197, '<':25137, '>':25137}
}

// Globals
var bix  map[rune]int // returns the bracket type ():0, []:1, {}:2, <>:3
var boc  map[rune]int // returns open (1) or close (-1)
var bpt  map[rune]int // returns the points for part 1 ):3, ]:57, }:1197, >:25137

// MAIN ----
func main () {

	start  := time.Now()
	ini()
	input  := readTxtFile("d10.input.txt")

	p1, p2 := navRepair(input)
	fmt.Printf("Sum of broken brackets: %v\n", p1)
	fmt.Printf("Auto-complete score: %v\n", p2)
 	fmt.Printf("Execution time: %v\n", time.Since(start))
}