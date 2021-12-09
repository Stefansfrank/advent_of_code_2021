package main

import (
	"fmt"
	"strconv"
	"os"
	"bufio"
	"time"
	//"regexp"
)

// no error handling ...
func readTxtFileInt (name string) (lines []int) {	
	file, _ := os.Open(name)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {	
		lines = append(lines, atoi(scanner.Text()))
	}
	return
}

// inline capable / no error handling Atoi
func atoi (x string) int {
	y, _ := strconv.Atoi(x)
	return y
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
	input  := readTxtFileInt("d1." + dataset + ".txt")

	// Part 1
	cnt    := 0
	for i:=1; i<len(input); i++ {
		if input[i] > input[i-1] {
			cnt++
		}
	}
	fmt.Printf("Increased %v times\n", cnt)

	// Part 2
	cnt     = 0
	rSum   := input[0]+input[1]+input[2]
	for i:=3; i<len(input); i++ {
		new := rSum - input[i-3] + input[i]
		if new > rSum {
			cnt++
		}
		rSum = new
	}
	fmt.Printf("Floating Average Increased %v times\n", cnt)

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}
