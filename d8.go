package main

import (
	"fmt"
	"math/bits"
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

// input parser
func parseFile (lines []string) (obsvs, disps [][]string) {

	re   := regexp.MustCompile(`([a-g]+)`)
    obsvs = [][]string{} // the 10 oberseved patterns for each input line (display)
    disps = [][]string{} // the four digits of the display

	for _, line := range lines {
        obs      := []string{}
        disp     := []string{}

        match    := re.FindAllString(line, -1)
        for i, m := range match {
            if i < 10 {
                obs  = append(obs, m)            
            } else {
                disp = append(disp, m)
            }
        }
		obsvs     = append(obsvs, obs)
        disps     = append(disps, disp)
	}
	return
}

// ---- problem independent helpers to handle the seven digit display --------------------------------

// This builds a bit representation of a seven digit display
// mapping the segments lit (as bits in a uint) to the number the pattern represents
// ... the only manually step 'explaining' the logic of the display
func build7D () (sevenD map[uint]int) {

    /*
      0:      1:      2:      3:      4:      5:      6:      7:      8:      9:
     aaaa    ....    aaaa    aaaa    ....    aaaa    aaaa    aaaa    aaaa    aaaa
    b    c  .    c  .    c  .    c  b    c  b    .  b    .  .    c  b    c  b    c
    b    c  .    c  .    c  .    c  b    c  b    .  b    .  .    c  b    c  b    c
     ....    ....    dddd    dddd    dddd    dddd    dddd    ....    dddd    dddd
    e    f  .    f  e    .  .    f  .    f  .    f  e    f  .    f  e    f  .    f
    e    f  .    f  e    .  .    f  .    f  .    f  e    f  .    f  e    f  .    f
     gggg    ....    gggg    gggg    ....    gggg    gggg    ....    gggg    gggg
      L6      L2      L5      L5      L4      L5      L6      L3      L7      L6            
    */

    // Map index: a uint where bit n represents the nth segment of the display (a = 0, b =1 ... in the above) 
    // Map value: the numeric digit represented by the segment pattern in the index

    sevenD = make(map[uint]int)
    sevenD[uint(0b01110111)] = 0
    sevenD[uint(0b00100100)] = 1
    sevenD[uint(0b01011101)] = 2
    sevenD[uint(0b01101101)] = 3
    sevenD[uint(0b00101110)] = 4
    sevenD[uint(0b01101011)] = 5
    sevenD[uint(0b01111011)] = 6
    sevenD[uint(0b00100101)] = 7
    sevenD[uint(0b01111111)] = 8
    sevenD[uint(0b01101111)] = 9
    return
}

// This uses the above sevenD to create a useful mapping array segCnt:
// For each of the 7 segments, it counts how many digits the segment is part of - subdivided by the total amount of segments lit for a given digit
// e.g. sevenD[0][5] = 2 means that the segment 0 ('aa' in the pic above) is part of 2 digits that require 5 segments to be lit
func buildSegCnt (svD map[uint]int) (segCnt [][]int) {

    // init
    segCnt = make([][]int, 7)
    for i := 0; i < 7; i++ {
        segCnt[i] = make([]int, 8)
    }

    // loops through the index of sevenD and counts how often each segment is used
    // the segment is represented as the bit position in this index
    for k, _  := range svD {
        bitLn := bits.OnesCount(k) // amount of signal wires that are lit for this digit

        // looping through the bits of the current digit 
        // and adding the counter for the length of the pattern if bit is set
        cBit := uint(0b00000001) 
        for i := 0; i < 7; i++ {
            if (cBit & k) > 0 {
                segCnt[i][bitLn] += 1
            }
            cBit <<= 1
        }
    }
    return
}

// computes the 'appearance pattern' of a wire  
// i.e. how often does this wire appear in observations of length [i]
// e.g. ap[5] = 2 means that this letter appears twice in signals with 5 segments lit
func appPat (wire byte, obs []string) (ap []int) {
    ap = make([]int, 8)
    for _, o := range obs {
        if strings.Contains(o, string(wire)) {
            ap[len(o)] +=1
        } 
    }
    return
}

// determine the wiring i.e. which letter is connected to which segment
// e.g. wiring['b'] = 0 means that the b-wire is connected to segment 0 ('aa' in the pic)
func detectWiring (obs []string, segCnt [][]int) (wiring map[byte]int) {
    wiring = make(map[byte]int)

    // loop over all letters
    for i := byte('a'); i < 'h'; i++ {
        ap := appPat(i, obs)

        // loop over all potential length counters
        for j,ss := range segCnt {

            // loop through the length counters that matter
            // and compare whether the typical pattern of segment j is the appearance pattern of wire i
            solution := true
            for k := 3; k<7; k++ {
                if ss[k] != ap[k] {
                    solution = false
                    break
                }
            }

            // found a match
            if solution {
                wiring[i] = j
                break
            }
        }     
    }
    return
} 

// displays the digit that is represented by a given wire signal (e.g. "acfg") using the given wiring
func digit(disp string, wiring map[byte]int, sevenD map[uint]int) int {

    var key uint // binary representation of lit segments

    // go through letters
    for _, c := range disp {
        key += 1 << wiring[byte(c)] // simple power of 2
    }

    return sevenD[key]
}

// MAIN ----
func main () {

    // helper for quick power of 10
    m10 := []int{1, 10, 100, 1000}

    // support of command line driven dataset picking
    dataset := ""
    if len(os.Args) < 2 || os.Args[1] == "" {
        fmt.Println("No argument given - trying 'test' dataset.")
        dataset = "test"
    } else {
        dataset = os.Args[1]
    }

	start  := time.Now()
	input  := readTxtFile("d8." + dataset + ".txt")
    obsvs, disps := parseFile(input) // the input: observations and displays

    // prepare the generic 7D helpers
    svD    := build7D()
	segCnt := buildSegCnt(svD)

    // Part 1 & 2 counter
    p1Cnt  := 0
    p2Cnt  := 0

    // loop through test sets
    for i, o := range obsvs {

        fmt.Printf("Disp %03v:   [", i)
        wiring := detectWiring(o, segCnt)
        for j := 0; j < 4; j++ {
            d := digit(disps[i][j], wiring, svD)
            if (d == 1 || d == 4 || d == 7 || d == 8) {
                p1Cnt += 1
            }
            p2Cnt += d * m10[3-j]
            fmt.Print(d)  
        }
        fmt.Printf("]   | P1 Result: %3v | P2 Result: %7v\n", p1Cnt, p2Cnt)
    }

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}