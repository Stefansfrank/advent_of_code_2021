package main

import (
	"fmt"
	"strconv"
	"os"
	"bufio"
	"time"
)

// no error handling ... 
// parses the binaries from the txt file straight into ints
func readTxtFileBInt (name string) (lines []int, nBit int) {	
	file, _ := os.Open(name)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		if nBit == 0 { // determine bit length on first non zero length input
			nBit = len(scanner.Text())
		}
		tmp, _ := strconv.ParseInt(scanner.Text(), 2, 32) // binary parsing		
		lines = append(lines, int(tmp))
	}
	return
}

// recursively cuts down the set of numbers while moving from the highest bit down
// return if only one number is in the set
// oxy: flag indiciating oxygen or CO2 rules
// ord: the position of the bit that should be counted for the most / least common detection
func getRating (nums []int, oxy bool, ord int) int {

	// end recursion if only one number is left
	if len(nums) == 1 {
		return nums[0]
	}

	// this clause exits if bit zero is reached and the end recursion condition has not been met
	// not strictly necessary as AoC does usually provide datasets leading to solutions ...
	if ord < 0 {
		return 0  
	}

	// build two stacks of numbers:
	// one for numbers with zeros and one for numbers with ones at the specified bit position
	zer := []int{}
	one := []int{}
	for _, num := range nums {
		if (num & (1 << ord)) == 0 {
			zer = append(zer, num)
		} else {
			one = append(one, num)
		}
	}

	// return the correct stack for the ruleset
	if (oxy) {
		if len(one) >= len(zer) {
			return getRating(one, oxy, ord-1)
		} else {
			return getRating(zer, oxy, ord-1)
		}
	} else {
		if len(zer) <= len(one) {
			return getRating(zer, oxy, ord-1)
		} else {
			return getRating(one, oxy, ord-1)
		}
	}

	// should never be reached but syntactical necessity
	return 0
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

	start     := time.Now()
	inp, nBit := readTxtFileBInt("d3." + dataset + ".txt")

	oxy    := getRating(inp, true,  nBit - 1)
	co2    := getRating(inp, false, nBit - 1)

 	fmt.Printf("Oxygen: %v, CO2: %v, Life support: %v\n", oxy, co2, oxy*co2)
 	fmt.Printf("Execution time: %v\n", time.Since(start))
}