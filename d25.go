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

// one step of moving sea cucumbers
func step(mp [][]byte) (changed bool) {

	// right moving
	for _, m := range mp {
		over := m[len(m) - 1] == '>' && m[0] == '.' // overflow?
		for ix := 0; ix < len(m) - 1; ix ++ {
			if m[ix] == '>' && m[ix+1] == '.' {
				m[ix+1] = '>'
				m[ix] 	= '.'
				ix += 1 // don't check the next spot 
				changed = true
			}
		}
		if over {
			m[len(m) - 1] = '.'
			m[0] = '>'
			changed = true
		}
	}

	// down moving
	for ix := 0; ix < len(mp[0]); ix++ {
		over := mp[len(mp) - 1][ix] == 'v' && mp[0][ix] == '.' // overflow?
		for iy := 0; iy < len(mp) - 1; iy++ {
			if mp[iy][ix] == 'v' && mp[iy+1][ix] == '.' {
				mp[iy+1][ix] = 'v'
				mp[iy][ix] 	= '.'
				iy += 1 // don't check the next spot
				changed = true
			}
		}
		if over {
			mp[len(mp) - 1][ix] = '.'
			mp[0][ix] = 'v'
			changed = true
		}
	}
	return
}

// print out map at the end to see whether something nice is hidden
func dump(mp [][]byte) {
	for _, l := range mp {
		for _, c := range l {
			fmt.Printf("%c",c)
		}
		fmt.Println()
	}
	fmt.Println()
}

// converts []string into [][]byte
func conv(inp []string) (mp [][]byte) {
	mp = make([][]byte, len(inp)) 
	for y := 0; y < len(inp); y++ {
		mp[y] = make([]byte, len(inp[y]))
		for x := 0; x < len(inp[y]); x++ {
			mp[y][x] = inp[y][x]
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
	input  := readTxtFile("d25." + dataset + ".txt")
	mp     := conv(input)

	var cnt int
	changed := true 
	for cnt = 0; changed; cnt++ {
		changed = step(mp)
	}

	fmt.Println("Stable after", cnt, "steps")
	dump(mp)

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}