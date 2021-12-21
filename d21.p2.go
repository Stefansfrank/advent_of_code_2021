package main

import (
	"fmt"
	"time"
)

var endScore  = 21
var boardSize = 10

// there are three core global variables that are [][][][]int 
// the inidices mean: [position player 1][position player 2][score player 1][score player 2]
// var stateCt counts how many different games are in the state described by the indices
// var diff  is a cache holding the changes in the amount of games in that state per round
var stateCt, diff [][][][]int 

// dist holds the amount of times a three roll total occurs for each end result
// winCnt counts the wins for each player
var dist, winCnt []int

// openStates counts the game states that have not yet ended
var openStates int

// initialize the game state counters
// ... sometimes I hate Go ...
func initial() {
	dim12 := boardSize + 1
	dim34 := endScore  + 10

	stateCt = make([][][][]int, dim12)
	diff    = make([][][][]int, dim12)
	for i := 0; i < 11; i++ {
		stateCt[i] = make([][][]int, dim12)
		diff[i]    = make([][][]int, dim12)
		for j :=0; j<11; j++ {
			stateCt[i][j] = make([][]int, dim34)
			diff[i][j]    = make([][]int, dim34)
			for k:=0; k<31; k++ {
				stateCt[i][j][k] = make([]int, dim34)
				diff[i][j][k]    = make([]int, dim34)
			}
		}
	}

	// the starting position
	stateCt[4][9][0][0] = 1
	openStates = 1
	winCnt = []int{0,0}

	// distribution of the results of the 27 possible roll combinations
	dist = make([]int,10)
	for i:=0; i<27; i++ {
		// using a 3 digit "trinary" representation of 27 to simulate the rolls
		dist[i/9 + i%9/3 + i%3 + 3] += 1 
	}
}

// executing one roll iteration with all possible outcomes
func movePlayer(plN int) {

	var cnt int
	dim12 := boardSize + 1
	dim34 := endScore
	oplN  := (plN + 1) % 2 // other player

	// these are used to hold the indices for access
	 pos := []int{0,0}
	 scr := []int{0,0}
	npos := []int{0,0}
	nscr := []int{0,0}

	// loop through all potential game states (of games with both scores < 21)
	for pos[0] = 1; pos[0] < dim12; pos[0]++ {
		for pos[1] = 1; pos[1] < dim12; pos[1]++ {
			for scr[0] = 0; scr[0] < dim34; scr[0]++ {
				for scr[1] = 0; scr[1] < dim34; scr[1]++ {
					cnt = stateCt[pos[0]][pos[1]][scr[0]][scr[1]]

					// look at game states only if there are more than 0 games in that state
					if cnt > 0  {

						// go through the potential roll results
						for r:=3; r < 10; r++ {

							// determine the new indices for the state after roll
							npos[plN]  = mov(pos[plN], r)
							npos[oplN] = pos[oplN]
							nscr[plN]  = scr[plN] + npos[plN]
							nscr[oplN] = scr[oplN]

							// cache the amount added to the new state							
							diff[npos[0]][npos[1]][nscr[0]][nscr[1]] += cnt * dist[r]

							// if that roll made a game finish
							if (nscr[plN]) > 20 {
								winCnt[plN] += cnt * dist[r]
							}
						}
					}

					// since these states have now evolved, we will need to subtract them
					diff[pos[0]][pos[1]][scr[0]][scr[1]] -= cnt
				}
			}
		}
	}

	// apply the cached summations
	var nzr bool
	for pos[0] = 1; pos[0] < dim12; pos[0]++ {
		for pos[1] = 1; pos[1] < dim12; pos[1]++ {
			for scr[0] = 0; scr[0] < dim34; scr[0]++ {
				for scr[1] = 0; scr[1] < dim34; scr[1]++ {

					// only deal with states that encountered a change
					nzr = diff[pos[0]][pos[1]][scr[0]][scr[1]] != 0
					if nzr {

						// if that state did not yet occur, add to the open state counter
						if stateCt[pos[0]][pos[1]][scr[0]][scr[1]] == 0 {
							openStates += 1
						}

						// apply changes (and reset them)
						stateCt[pos[0]][pos[1]][scr[0]][scr[1]] += diff[pos[0]][pos[1]][scr[0]][scr[1]]
						diff[pos[0]][pos[1]][scr[0]][scr[1]] = 0

						// if that state no longer occurs, remove them from the open state counter
						if stateCt[pos[0]][pos[1]][scr[0]][scr[1]] == 0 {
							openStates -= 1
						}

					}
				}
			}
		}
	}				
}

// quick helper adding two numbers on the cyclic 1..10 board
func mov(from, by int) int {
	return (from + by - 1) % 10 + 1
}

func main() {

	start := time.Now() 

	initial()
	for openStates > 0 {
		movePlayer(0)
		movePlayer(1)
	}

	fmt.Println("Player 1 wins", winCnt[0], "times, Player 2 wins", winCnt[1], "times")
	fmt.Println("Execution time:", time.Since(start))

}