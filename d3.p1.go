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

// calculate Gamma & Epsilon
func getGamEps (lines []string) (gam, eps int) {
	nBit := len(lines[0])
	nLns := len(lines)
	cnt  := make([]int, nBit) // a counter for each bit

	// adding 1 to the counter of the relevant bit for each '1' in each number 
	for _, line := range lines {
		for i := 0; i < nBit; i++ {
			cnt[i] += int(byte(line[i])) - 48 // ascii of '0'
		}
	}

	// determining the result by the observation that if the counter for a given bit
	// is more than half of the total number of input numbers, the prevalent bit is '1'
	gam  = 0
	msk := 0 // for the inversion at the end, I need an integer with '1' at all positions used by the input
	for i := 0; i < nBit; i++ {
		bas := 1 << (nBit-i-1)
		if cnt[i] > nLns/2 {
			gam += bas // manual binary -> decimal conversion as I am looping through bits anyway
		}
		msk += bas // ... all relevant bits are '1' in the mask
	}
	eps = msk ^ gam // epsilon is by definition the inverse of gamma

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
	input  := readTxtFile("d3." + dataset + ".txt")
	gam, eps := getGamEps(input)

 	fmt.Printf("Gamma: %v, Epsilon: %v, Power:%v\n", gam, eps, gam*eps)
 	fmt.Printf("Execution time: %v\n", time.Since(start))
}