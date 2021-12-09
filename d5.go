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

// simple point structure (index for the map)
type pnt struct {
	x int
	y int
}

// builds map from the input
func buildMap(lines []string, p2 bool) (result map[pnt]int) {

	re  := regexp.MustCompile(`([0-9]+),([0-9]+) -> ([0-9]+),([0-9]+)`)
	mp  := make(map[pnt]int)

	for _, line := range lines {
		match := re.FindStringSubmatch(line)
		xf   := atoi(match[1])
		yf   := atoi(match[2])
		xt   := atoi(match[3])
		yt   := atoi(match[4])

		// kick out diagonals for part 1
		if !p2 && xf != xt && yf != yt {
			continue
		}

		// determine incrementors for both axes
		xinc := 0
		yinc := 0
		if xf > xt {
			xinc = -1
		} else if xf < xt {
			xinc = 1
		}
		if yf > yt {
			yinc = -1
		} else if yf < yt {
			yinc = 1
		}

		// go along the line
		for (true) {
			mp[pnt{x: xf, y: yf}] += 1
			if xf == xt && yf == yt {
				break
			}
			xf += xinc
			yf += yinc
		}
	}
	return mp
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
	input  := readTxtFile("d5." + dataset + ".txt")

	// part 1
	mp  := buildMap(input, false)
	cnt := 0
	for _, v := range mp {
		if v >= 2 {
			cnt += 1
		}
	}
	fmt.Printf("Two or more horizontal / vertical lines cross %v times\n", cnt)

	// part 2
	mp  = buildMap(input, true)
	cnt = 0
	for _, v := range mp {
		if v >= 2 {
			cnt += 1
		}
	}
	fmt.Printf("Two or more total lines cross %v times\n", cnt)
 	fmt.Printf("Execution time: %v\n", time.Since(start))
}