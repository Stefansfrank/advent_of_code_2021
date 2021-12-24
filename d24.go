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

// AOI simulator
func run(code []string, inp []int) (reg map[byte]int) {

	reg  = map[byte]int{}
	inx := 0

	for i,line := range code {
		cmd := line[:3]
		tgt := line[4]
		opr := ""
		num := false
		if len(line) > 6 {
			opr = line[6:]
			num = line[6] < byte(100)
		}

		switch cmd{
		case "inp":
			reg[tgt] = inp[inx]
			inx += 1
		case "add":
			if num {
				reg[tgt] += atoi(opr)
			} else {
				reg[tgt] += reg[opr[0]]
			}
		case "mul":
			if num {
				reg[tgt] *= atoi(opr)
			} else {
				reg[tgt] *= reg[opr[0]]
			}
		case "div":
			if num {
				reg[tgt] = reg[tgt] / atoi(opr)
			} else {
				reg[tgt] = reg[tgt] / reg[opr[0]]
			}
		case "mod":
			if num {
				reg[tgt] = reg[tgt] % atoi(opr)
			} else {
				reg[tgt] = reg[tgt] % reg[opr[0]]
			}
		case "eql":
			var eq bool
			if num {
				eq = reg[tgt] == atoi(opr)
			} else {
				eq = reg[tgt] == reg[opr[0]]
			}
			if eq {
				reg[tgt] = 1
			} else {
				reg[tgt] = 0
			}
		}
		if (i + 1) % 18 == 0 {
			fmt.Println("Regs after line", i, "-", reg)
		}
	}
	return
}

// analyzes the input on variations between the handling of the 14 inputs
// using the fact that for each inp there are 17 very similar commands after
func analyze(code []string) {

	for ix := 0; ix < 18; ix++ {

		cmd := make([]string, 14)
		tgt := make([]string, 14)
		opr := make([]string, 14)

		// parse command, target and operand
		for it := 0; it < 14; it++ {
			cmd[it] = code[it * 18 + ix][:3]
			tgt[it] = code[it * 18 + ix][4:5]
			if len(code[it * 18 + ix]) > 6 {
				opr[it] = code[it * 18 + ix][6:]
			}
		}

		// analyze if all are identical
		ccmd := cmd[0]
		ttgt := tgt[0]
		oopr := opr[0]
		oneCmd := true
		oneTgt := true
		oneOpr := true
		for it := 1; it <14; it++ {
			if cmd[it] != ccmd {
				oneCmd = false
			}
			if tgt[it] != ttgt {
				oneTgt = false
			}
			if opr[it] != oopr {
				oneOpr = false
			}
		}
		if oneCmd {
			cmd = []string{ccmd}
		}
		if oneTgt {
			tgt = []string{ttgt}
		}
		if oneOpr {
			opr = []string{oopr}
		}

		fmt.Println(cmd, tgt, opr)
	}	
}

// initializes the input with 14 identical digits n
func initial(n int) (is []int) {
	is = make([]int, 14)
	for i := 0; i < 14; i++ {
		is[i] = n
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
	code   := readTxtFile("d24." + dataset + ".txt")

	analyze(code)
	fmt.Println()

	for i := 1; i < 10; i++ {
		fmt.Printf("All %vs ------------------\n", i)
		run(code, initial(i))		
	}

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}