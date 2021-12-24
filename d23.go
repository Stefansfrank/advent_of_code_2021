package main

import (
	"fmt"
	"time"
)

type xy struct {
	x int
	y int
}

func initial(tst bool, anum int) (amp []xy, loc map[xy]int, en, gol []int, rmY map[xy]bool) {

// #############
// #01234567890#
// ###1#1#1#1###
//   #2#2#2#2#
//   #########

	if tst {
		if anum == 2 {
			amp = []xy{ xy{2,2}, xy{8,2},
						xy{2,1}, xy{6,1},
						xy{4,1}, xy{6,2},
						xy{4,2}, xy{8,1} }
		} else {
			amp = []xy{ xy{2,4}, xy{8,4}, xy{8,2}, xy{6,3},
						xy{2,1}, xy{6,1}, xy{6,2}, xy{4,3},
						xy{4,1}, xy{6,4}, xy{4,2}, xy{8,3},
						xy{4,4}, xy{8,1}, xy{2,2}, xy{2,3} }
		}
	} else {
		if anum == 2 {
			amp = []xy{ xy{4,1}, xy{4,2},
						xy{6,1}, xy{2,2},
						xy{8,1}, xy{8,2},
						xy{2,1}, xy{6,2} }
		} else {
			amp = []xy{ xy{4,1}, xy{4,4}, xy{8,2}, xy{6,3},
						xy{6,1}, xy{2,4}, xy{6,2}, xy{4,3},
						xy{8,1}, xy{8,4}, xy{4,2}, xy{8,3},
						xy{2,1}, xy{6,4}, xy{2,2}, xy{2,3} }
		}
	}

	loc = make(map[xy]int)
	for i, a := range amp {
		loc[a] = i
	}
	for i := 0; i < 11; i++ {
		loc[xy{i, 0}] = -1
	}

	// energy values
	mult := 1
	en  = make([]int, 0, 4*anum)
	for i := 0; i < 4; i++ {
		for j := 0; j < anum; j++ {
			en = append(en, mult)

		}
		mult *= 10
	}

	// these are the goal rooms
	gol = make([]int, len(amp))
	for i := 0; i < len(amp); i++ {
		gol[i] = 2 + 2 * (i/anum)
	}

	rmY = map[xy]bool{xy{2,0}:true, xy{4,0}:true, xy{6,0}:true, xy{8,0}:true}

	return
}

type state struct {
	amp []xy       // a list of the coordinates of each amphi
	loc map[xy]int // whats on the map at location xy (redundant with the above but faster)
	ast []int  	   // state of amphipod { 0: start position, 1: hallway, 2: goal}
	eng int 	   // total energy
	typ int    	   // type of state { -1: failed, 0: ongoing, 1: finished}
	trc string	   // a string representation of what happened to get at this state (debugging only)
}

// create a new copy of a slice of type int
func cpint(isl []int) (nisl []int) {
	nisl = make([]int, len(isl))
	for ix, in := range(isl) {
		nisl[ix] = in
	}
	return
}

// create a new copy of a slice of type map[xy]int
func cpmap(mpp map[xy]int) (nmpp map[xy]int) {
	nmpp = make(map[xy]int)
	for k,v := range(mpp) {
		nmpp[k] = v
	}
	return
}

// create a new copy of a slice of type []xy
func cpxy(xyl []xy) (nxyl []xy) {
	nxyl = make([]xy, len(xyl))
	for ix, xy := range(xyl) {
		nxyl[ix] = xy
	}
	return
}

// only for debugging - prints the state with all information
func (st state) dump(anum int) {
	fmt.Println("*************")
	fmt.Print("*")
	for i := 0; i < 11; i++ {
		ll := st.loc[xy{i, 0}]
		if ll == -1 {
			fmt.Print(".")
		} else {
			fmt.Printf("%c", byte(ll / anum) + 'A')
		}
	}
	fmt.Println("*")
	for i := 1; i <= anum; i ++ {
		fmt.Print("***")
		ll := st.loc[xy{2, i}]
		if ll == -1 {
			fmt.Print(".")
		} else {
			fmt.Printf("%c", byte(ll / anum) + 'A')
		}
		fmt.Print("*")
		ll  = st.loc[xy{4, i}]
		if ll == -1 {
			fmt.Print(".")
		} else {
			fmt.Printf("%c", byte(ll / anum) + 'A')
		}
		fmt.Print("*")
		ll  = st.loc[xy{6, i}]
		if ll == -1 {
			fmt.Print(".")
		} else {
			fmt.Printf("%c", byte(ll / anum) + 'A')
		}
		fmt.Print("*")
		ll  = st.loc[xy{8, i}]
		if ll == -1 {
			fmt.Print(".")
		} else {
			fmt.Printf("%c", byte(ll / anum) + 'A')
		}
		fmt.Println("***")
	}
	fmt.Println("*************")
	fmt.Println("Amp", st.amp)
	fmt.Println("Ast", st.ast)
	fmt.Println("Typ",st.typ,"Eng",st.eng)
	fmt.Println("Trace", st.trc)
}

// creates a unique string hash for each state
func (st state) hash() (h string) {
	for _, a := range st.amp {
		h += fmt.Sprintf("[%v,%v]",a.x,a.y) 
	}
	return
}

// the central function
func simulate(test bool, anum int) {

	// initializes the start data
	iamp, iloc, en, gol, rmY := initial(test, anum)
	states := []state{state{amp:iamp, loc:iloc, ast: make([]int, len(iamp)), typ:int(0)}}

	// handles test data that have amphis already in their final position at the start
	for i,a := range states[0].amp {
		if a.x == 2 + 2*(i/2) && a.y == anum {
			states[0].ast[i] = 2
		}
	}

	wins   := []state{}        // collects the winning states
	cach   := map[string]int{} // Maps a unique hash for the state to the energy needed to get into that state
	var nstates []state        // swapping list as the state list will be recreated on each run

	for cnt := 0; len(states) > 0 || cnt < 30; cnt++ {
		fmt.Println("Iteration", cnt, "starts with", len(states), "states")

		// go through all states
		nstates = []state{}
		ns     := state{}
		for stix,st := range states {

			// the number of states so I can see whether I added new ones
			nmst := len(nstates)

			// if there happen to be non open states here 
			if st.typ != 0 {
				continue
			}

			// through all amphis in a state
			aloop: for ai, a := range st.amp {

				// don't do anything if the amphi is settled
				if st.ast[ai] != 2 {

					nen := 0 // new energy

					// still in start room
					if st.ast[ai] == 0 {

						// can I get out?
						for y := a.y - 1; y > 0; y-- {
							if st.loc[xy{a.x, y}] != -1 {
								//fmt.Println("Can't get out", ai, a.x, a.y, y)
								continue aloop
							}
						} 

						// yes I can get out, add the cost of moving up into the hall
						nen += a.y * en[ai]

						// look for another hall position to the left
						for x := a.x; x > -1; x-- {
							tp := xy{x, 0}
							if st.loc[tp] != -1 {
								break
							} else if rmY[tp] {
								continue
							} else {

								ns = state{amp:cpxy(st.amp), loc:cpmap(st.loc), ast:cpint(st.ast), eng:st.eng, typ:st.typ}
								ns.eng += nen + en[ai] * (a.x - x) 
								ns.loc[a] = -1
								ns.amp[ai] = tp 
								ns.loc[tp] = ai
								ns.ast[ai] = 1
								ns.trc = st.trc + fmt.Sprintf("[%v] to [%v,%v](%v)|", ai, tp.x, tp.y, ns.eng)

								oeng, fnd := cach[ns.hash()]
								if !fnd || oeng > ns.eng {
									cach[ns.hash()] = ns.eng
									nstates = append(nstates, ns)
								} 
							}
						}

						// look for another hall position on the right
						for x := a.x; x < 11; x++ {
							tp := xy{x, 0}
							if st.loc[tp] != -1 {
								break
							} else if rmY[tp] {
								continue
							} else {

								ns = state{amp:cpxy(st.amp), loc:cpmap(st.loc), ast:cpint(st.ast), eng:st.eng, typ:st.typ}
								ns.eng += nen + en[ai] * (x - a.x) 
								ns.loc[a] = -1
								ns.amp[ai] = tp 
								ns.loc[tp] = ai
								ns.ast[ai] = 1
								ns.trc = st.trc + fmt.Sprintf("[%v] to [%v,%v](%v)|", ai, tp.x, tp.y, ns.eng)

								oeng, fnd := cach[ns.hash()]
								if !fnd || oeng > ns.eng {
									cach[ns.hash()] = ns.eng
									nstates = append(nstates, ns)
								} 
							}
						}
					}

					// now see whether I can get this into the goal room
					tx := true

					// determines direction and catches a case where the lower amphi moves up
					stp := 1
					if gol[ai] < a.x {
						stp = -1
					} else if gol[ai] == a.x {
						tx = false
					}

					// can I get to the room?
					for x := a.x + stp; x != gol[ai] && tx; x += stp {
						if st.loc[xy{x, 0}] != -1 {
							tx = false
						}
					}

					// I can get in front of the room
					if tx {

						// now see whether the room is open for me
						ty := 0
						for gy := 1; gy <= anum; gy++ {
							if st.loc[xy{gol[ai], gy}] != -1 && st.loc[xy{gol[ai], gy}]/anum != ai/anum {
								ty = 0
								break
							}
							if st.loc[xy{gol[ai], gy}] == -1 {
								ty = gy
							}
						}

						// open !!
						if ty > 0 {
							ns = state{amp:cpxy(st.amp), loc:cpmap(st.loc), ast:cpint(st.ast), eng:st.eng, typ:st.typ}
							ns.eng += nen + en[ai] * (gol[ai] - a.x) * stp + ty * en[ai]
							ns.loc[a] = -1
							ns.amp[ai] = xy{gol[ai], ty} 
							ns.loc[ns.amp[ai]] = ai
							ns.ast[ai] = 2
							ns.trc = st.trc + fmt.Sprintf("[%v] to [%v,%v](%v)|", ai, gol[ai], ty, ns.eng)

							// win detection
							tmp := 0
							for _,s := range ns.ast {
								tmp += s
							}
							if tmp >= 2*len(ns.ast) {
								ns.typ = 1
								wins = append(wins, ns)
							}
							
							oeng, fnd := cach[ns.hash()]
							if !fnd || oeng > ns.eng {
								cach[ns.hash()] = ns.eng
								nstates = append(nstates, ns)
							} 
						}
					}
				}
			}

			// found a stuck state
			if len(nstates) == nmst {
				states[stix].typ = -1
			}
			
		}
		states = nstates
	}

	min := 1000000000
	for _, s := range wins {
		if s.eng < min {
			min = s.eng
		}
	}
	fmt.Println("Min:", min)
}

func main() {

	start := time.Now()

	simulate(true, 4)

	fmt.Println("Exeution time:", time.Since(start))
}
