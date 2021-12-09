package main

import (
	"fmt"
	"strconv"
	"os"
	"bufio"
	"time"
	"regexp"
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

// no error handling ...
func readCsvTxtFileInt (name string) (nums [][]int) {	
	nums := [][]int{}
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

// input parser using Regex
func parseFile (lines []string) (result int) {

	re  := regexp.MustCompile(``)

	for _, line := range lines {
		match := re.FindStringSubmatch(line)
		//parse
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

	start  := time.Now()
	input  := readTxtFile("d." + dataset + ".txt")
	input  := readTxtFileInt("d." + dataset + ".txt")
	input  := readCsvTxtFileInt("d." + dataset + ".txt")
	result := parseFile(input)

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}