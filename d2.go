package main

import (
	"fmt"
	"strconv"
	"os"
	"bufio"
	"time"
	"regexp"
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

// execute commands
func execute (lines []string, p2 bool) (pos, dep int) {

	re  := regexp.MustCompile(`(forward|down|up) ([0-9]+)`)
	aim := 0

	for _, line := range lines {
		match := re.FindStringSubmatch(line)
		switch match[1] {
		case "forward":
			del := atoi(match[2])
			pos += del
			if p2 {
				dep += aim * del
			}
		case "up":
			if p2 {
				aim -= atoi(match[2])
			} else {
				dep -= atoi(match[2])
			}
		case "down":
			if p2 {
				aim += atoi(match[2])
			} else {
				dep += atoi(match[2])
			}
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

	start  := time.Now()
	input  := readTxtFile("d2." + dataset + ".txt")

	pos, dep := execute(input, false)
	fmt.Printf("Final position (P1 Rules): %v, depth: %d. Result: %v\n",pos,dep,pos*dep)

	pos, dep  = execute(input, true)
	fmt.Printf("Final position (P2 Rules): %v, depth: %d. Result: %v\n",pos,dep,pos*dep)

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}