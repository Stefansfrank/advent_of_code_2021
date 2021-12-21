package main

import (
	"fmt"
	"time"
)

const maxDice  = 1000 // when does the deterministic dice roll over
const ringSize = 10   // the size of the board
const rollsPer = 3    // how often is the dice rolled
const endScore = 1000 // the ending score
const numPlayr = 2    // the number of players

// the dice 
type detDc struct {
	state int
}

// rolling the dice 'rollsPer' times
func (d *detDc) roll() (res int) {
	for i := 0; i < rollsPer; i++ {
		d.state = (d.state + 1) % maxDice
		res += d.state	
	}
	return 
}

// the player
type player struct {
	pos int
	scr int
}

// move the player by rolling the dice
func (p *player) move(f int) {
	p.pos = (p.pos + f - 1) % ringSize + 1
	p.scr += p.pos
}

func main() {
	
	start := time.Now()

	// starting positions (4,9 was my input)
	pl := []*player{ &player{pos:4}, &player{pos:9} }
	d  := &detDc{state:0}

	rolls := 0

	// game moves
	outer: for {

		// player loop
		for i := 0; i < numPlayr; i++ {
			
			pl[i].move(d.roll())
			rolls += rollsPer
			if pl[i].scr >= endScore {
				lsr := (i + 1) % 2
				fmt.Printf("Game ends after %v rolls - puzzle output: %v (looser score:%v)\n",
							 rolls, rolls * pl[lsr].scr, pl[lsr].scr)
				break outer
			}		
		}
	}

	fmt.Println("Execution time:", time.Since(start))
}