package main

import (
	"fmt"
	"strconv"
	"os"
	"bufio"
	"time"
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

// inline capable / no error handling Atoi
func atoi (x string) int {
	y, _ := strconv.Atoi(x)
	return y
}

// input parser creating pGraph
// pGraph maps a cave to all possible caves to be reached from there
// pGraph["abc"] returns a slice with all connected cave names
func parseFile (lines []string) {

	pGraph = make(map[string][]string, len(lines))

	for _, line := range lines {
		ix := strings.Index(line, "-")
		pGraph[line[:ix]]   = append(pGraph[line[:ix]], line[ix+1:]) 
		pGraph[line[ix+1:]] = append(pGraph[line[ix+1:]], line[:ix]) 
	}
}

// makes a real copy of the "visited" map that is part of each sequence of paths keeping
// track of which caves are already visited
func copyVisPlus(vis map[string]bool, des string) (nVis map[string]bool) {
	nVis = make(map[string]bool)
	for k, v := range vis {
		nVis[k] = v
	}
	nVis[des] = true
	return
}

// global types / variables

// Represents a sequence of caves
type seq struct {
	len int      // total length  
	des string   // the final destination ("end" for valid ones)
	dmp string   // a string representation of the sequence (for debugging)
	vis map[string]bool   // marks all caves visited in this sequence
	ext bool     // if this is false an extra visit to one small cave is still possible
}
var pGraph map[string][]string 
var sGraph []seq // list of all valid sequences

// This is the core of the logic where all possible sequences of cave visits are build
// the logic is depth-first i.e. one sequence is build until "end" or no more possible next paths
// any branches encountered on the way are added at the end of the sequence list in order to
// be followed later. The list of all sequences is local (var tGr []seq) while only paths reaching "end"
// are copied over to the global (var sGraph []seq)
func buildSGraph(part1 bool) {
	tGr   := []seq{}
	sGraph = []seq{}

	// the initial caves reachable from start are added to the potential sequence list
	for _, des := range pGraph["start"] {
		tGr = append(tGr, seq{len: 1, des: des, dmp: "start-" + des, 
							vis: map[string]bool{"start": true, des: true}, ext: part1})
	}

	currSq := 0
	// this loop tries to add a step to the current sequence, adds branched sequences for later 
	// and goes to the next sequence whenever the current is terminated (ended or stuck)
	for (currSq < len(tGr)) {
		sq := tGr[currSq] // just short for read access to the current path

		// determine "end" to detect a valid sequence
		if sq.des == "end" {
			sGraph = append(sGraph, sq)
			currSq += 1
			continue
		}

		// this identifies valid next steps
		next := []string{}
		for _, des := range pGraph[sq.des] {

			// This catches all invalid next caves that have a connection. Invalid are
			// small letter caves that are already visited and there is no more extra visit or they are "start"
			if (rune(des[0]) > 'Z' && sq.vis[des] && (sq.ext || des == "start")) { 
				continue
			}

			// otherwise, this is a valid next step
			next = append(next, des)
		}

		// if more than one next step, add the potential other next steps  
		// to the current sequence and add each to the end of the potential sequence
		if len(next) > 1 {
			for i := 1; i < len(next); i++ {
				tGr   = append(tGr, seq{len: sq.len+1, des: next[i],
					dmp: sq.dmp + "-" + next[i], vis: copyVisPlus(sq.vis, next[i]),
					ext: sq.ext || (rune(next[i][0]) > 'Z' && sq.vis[next[i]])})
			}
		}

		// continue the current sequence with the 1st identified next step
		if len(next) > 0 {
			tGr[currSq].len = sq.len+1
			tGr[currSq].des = next[0]
			tGr[currSq].dmp = sq.dmp + "-" + next[0]
			tGr[currSq].ext = sq.ext || (rune(next[0][0]) > 'Z' && sq.vis[next[0]])
			tGr[currSq].vis[next[0]] = true

		// you are stuck, go look at next sequence
		} else {
			currSq += 1
		}
	}
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
	input  := readTxtFile("d12." + dataset + ".txt")
	parseFile(input)

	buildSGraph(true)
	fmt.Println("sGraph has", len(sGraph),"possible sequences." )

	buildSGraph(false)
	fmt.Println("sGraph has", len(sGraph),"possible sequences." )

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}