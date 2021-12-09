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

// input parser creates the school
// the scool is represented by 9 counters representing the amount of fish at each age
func parseFile (lines []string) (school []int) {

	school = make([]int, 9)
	slcd  := strings.Split(lines[0],",")
	for _, fs := range slcd {
		school[atoi(fs)] += 1
	}
	return
}

// counts the fish in a school
func countFish (school []int) (cnt int) {
	for _, fs := range school {
		cnt += fs
	}
	return
}

// iterates a school by one day
func nextDay (school []int) []int {

	nw := school[0] // amount of new fish created (= amount of fish age 0)

	// reduces the age of each fish be movin the down to the next lower counter
	// - the age 0 fish are cached in 'nw' so they will be added later
	// - the age 8 fish will be replaced by the newborns later
	for i := 1; i < 9; i++ {
		school[i-1] = school[i]
	}

	school[8]  = nw // newly created fish (overwriting existing age 8 population)
	school[6] += nw // new fish parents go back and are added to age 6 counter

	return school
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
	input  := readTxtFile("d6." + dataset + ".txt")
	school := parseFile(input)

	for i := 0; i < 256; i++ {
		school = nextDay(school)	
		if i == 79 || i == 255 {
			fmt.Printf("Size of school at day %v: %v fish.\n", i+1, countFish(school))
		}
	}

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}