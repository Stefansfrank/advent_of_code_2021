package main

import (
	"fmt"
	"strconv"
	"os"
	"bufio"
	"time"
	"strings"
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

// representation of a board
type board struct {
	num  [][]int // the numbers on the board
	hit  [][]int // 1 if the number was hit
	bng  int     // the round this board went bingo
	scr  int     // the score when the board went bingo
}

// creates a new board
func newBoard() board {
	ret := board{}
	ret.num = make([][]int,5)  
	ret.hit = make([][]int, 5)
	for i := 0; i < 5; i++ {
		ret.num[i] = make([]int,5)
		ret.hit[i] = make([]int,5)
	}
	ret.bng = 0
	ret.scr = 0
	return ret
}

// input parser building the draw sequence and the boards
func parseFile (lines []string) {

	// draw sequence
	drawS := strings.Split(lines[0], ",")
	draw = []int{}
	for _, ds := range drawS {
		draw = append(draw, atoi(ds))
	}

	// boards
	re  := regexp.MustCompile(` ?([0-9]+)  ?([0-9]+)  ?([0-9]+)  ?([0-9]+)  ?([0-9]+)`)
	boards = []board{}

	lc := 1 // running line counter (indexing the empty line before the board)
	for lc < len(lines) - 1 {
		nb     := len(boards)
		boards  = append(boards, newBoard())
		for i := 1; i < 6; i++ {
			line  := lines[lc+i]
			match := re.FindStringSubmatch(line)
			for j := 0; j < 5; j++ {
				boards[nb].num[i-1][j] = atoi(match[j+1])
			}
		}
		lc += 6
	}
}

// detects a bingo and sets the round indicator and the score at time of bingo
func (b *board) bingo(dIx int) {
	for i := 0; i < 5; i++ {
		if b.hit[i][0] + b.hit[i][1] + b.hit[i][2] + b.hit[i][3] + b.hit[i][4] == 5 {
			b.bng = dIx
			b.score(dIx)
		}
		if b.hit[0][i] + b.hit[1][i] + b.hit[2][i] + b.hit[3][i] + b.hit[4][i] == 5 {
			b.bng = dIx
			b.score(dIx)
		}
	}
}

// detects a hit for draw dIx and sets the hit counter
func (b *board) addDraw(dIx int) {
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			if b.num[i][j] == draw[dIx] {
				b.hit[i][j] = 1
				b.bingo(dIx)
			}
		}
	}
}

// computes the score at a given draw
func (b *board) score(dIx int) {
	b.scr = 0
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			if b.hit[i][j] == 0 {
				b.scr += b.num[i][j]
			}
		}
	}
	b.scr *= draw[dIx]
}

// the draw sequence and the boards
var draw   []int
var boards []board

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
	input  := readTxtFile("d4." + dataset + ".txt")
	parseFile(input)

	for dIx, _ := range draw {
		for bIx, _ := range boards {
			if boards[bIx].bng == 0 {
				boards[bIx].addDraw(dIx)
			}
		}
	}

	lastScore   := 0
	lastBngRnd  := 0
	firstScore  := 0
	firstBngRnd := len(draw)+1

	for _ ,brd := range boards {
		if brd.bng > lastBngRnd {
			lastScore  = brd.scr 
			lastBngRnd = brd.bng
		}
		if brd.bng < firstBngRnd {
			firstScore  = brd.scr 
			firstBngRnd = brd.bng
		}
	}

	fmt.Printf("First board scored %v, last board scored %v\n", firstScore, lastScore)
 	fmt.Printf("Execution time: %v\n", time.Since(start))
}