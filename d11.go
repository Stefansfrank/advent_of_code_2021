package main

import (
	"fmt"
	"strconv"
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

// inline capable / no error handling Atoi
func atoi (x string) int {
	y, _ := strconv.Atoi(x)
	return y
}

// input parser creating octopedes energy matrix
func parseFile (lines []string) (m intM) {

	m = newIntM(len(lines[0]), len(lines))
	for y, line := range lines {
		for  x, dig := range line {
			m.set(x, y, atoi(string(dig)))
		}
	}
	return
}

// basic matrix helpers ------------------------

// a new matrix of type int
func newIntM (xd, yd int) intM {
	m := make([][]int, yd)
	for y := 0; y< yd; y++ {
		m[y] = make([]int, xd)
	}
	return m
}

// a new mask (boolean matrix) for a given matrix
func newMask (m intM) boolM {
	xd := len(m[0])
	yd := len(m)
	nm := make([][]bool, yd)
	for y := 0; y< yd; y++ {
		nm[y] = make([]bool, xd)
	}
	return nm
}

// setters and getters in order to keep the order of the indices clean
func (m intM) set(x, y, v int) {
	m[y][x] = v
}
func (m boolM) set(x, y int, v bool) {
	m[y][x] = v
}
func (m intM) get(x, y int) int {
	return m[y][x]
}
func (m boolM) get(x, y int) bool {
	return m[y][x]
}

// simple add function
func (m intM) add(x, y, v int) {
	m[y][x] += v
}

// prints out the integer matrix
func (m intM) dump() {
	for _,line := range m {
		for _, dig := range line {
			fmt.Printf("%v", dig)
		}
		fmt.Println()
	}
}

// prints out the boolean matrix
func (m boolM) dump() {
	for _,line := range m {
		for _, dig := range line {
			if dig {
				fmt.Print("1")
			} else {
				fmt.Print("0")
			}
		}
		fmt.Println()
	}
}

// pads the matrix with value v in all directions
// returns a new matrix
func (m intM) padInt(v int) (nm intM) {
	xd := len(m[0])
	yd := len(m)

	nm       = make([][]int, yd+2)
	nm[0]    = make([]int, xd+2)
	nm[yd+1] = make([]int, xd+2)
	for i := 0; i < (xd+2); i++ {
		nm[0][i]    = v
		nm[yd+1][i] = v		
	}
	for i := 1; i < (yd+1); i++ {
		nm[i] = append([]int{v}, m[i-1]...)	
		nm[i] = append(nm[i], v)	
	}
	return nm
}

type intM  [][]int
type boolM [][]bool

// core logic
// either executes k steps and returns total octopedes lit across all steps (part 1)
// or (if k = -1) executes until all octopedes are lit and returns step count (part 2)
func steps(octs intM, k int) int {

	xd   := len(octs[0]) // dimensions of unpadded matrix
	yd   := len(octs)
	m    := octs.padInt(0)
	cnt  := 0
	lit  := 0

	// steps
	for i := 0; (k == -1) || (i < k); i++ {

		mask     := newMask(m) // mask is true if octopus lit in this step
		if (k == -1) { 
			lit = 0 // part 2
		}

		// initial increase by 1 for ever octopus
		for x:=1; x<=xd; x++ {
			for y:= 1; y<=yd; y++ {
				m.set(x,y,m.get(x, y) + 1)
			}
		}

		// action detects whether an octops was lit thus necessitating another 
		// loop scanning for secondary lits (initially true to force first scan)
		action := true
		for action {
			action = false

			// light up octopedes
			for x:=1; x<=xd; x++ {
				for y:= 1; y<=yd; y++ {
					if m.get(x,y) > 9 && !mask.get(x,y) {
						lit += 1
						m.add(x+1,y,1)
						m.add(x-1,y,1)
						m.add(x,y+1,1)
						m.add(x,y-1,1)
						m.add(x+1,y+1,1)
						m.add(x-1,y-1,1)
						m.add(x-1,y+1,1)
						m.add(x+1,y-1,1)
						mask.set(x,y,true)
						action = true
					}
				}
			}

			// clean up by forcing all lit octobedes to zero
			for x:=1; x<=xd; x++ {
				for y:= 1; y<=yd; y++ {
					if mask.get(x,y) {
						m.set(x,y,0)
					} 
				}
			}
		}

		cnt += 1
		if k == -1 && lit == xd*yd {
			return cnt // part 2
		}
	}
	return lit // part 1
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
	input  := readTxtFile("d11." + dataset + ".txt")
	octs   := parseFile(input)

	// Part 1
	stps   := 100
	fmt.Printf("Octopedes lit after %v steps: %v\n", stps, steps(octs, stps))

	// Part 2
	fmt.Printf("All octopedes lit in step %v\n", steps(octs, -1))

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}