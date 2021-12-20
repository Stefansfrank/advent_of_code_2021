package main

import (
	"fmt"
	"strconv"
	"os"
	"bufio"
	"time"
	"sort"
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

// types of pairs
const noP    = 0 // no pointers, both values
const rightP = 1 // the left is a value, the right a pointer to a pair
const leftP  = 2 // the left is a pointer to a pair, the right a value
const bothP  = 3 // both elements of the pair are pointers

// core pair structure
type pair struct {
	typ  int     // Type of pair - indicates what the two elements are, pointers to another pair or values (see constants below for values) 
	lvl  int     // the level / distance from the root pair which is lvl = 0
	val  []int   // the two values (size 2, will not be both filled unless type 0/noP)
	prs  []*pair // the two pointers to pairs (size 2 will not be both fulled unless type 3/bothP)
	up   *pair   // points to the par hierachically up from the pair (nil for root pair)
	upIx int     // indicates whether this pair is linked as the left or right element in the parent pair (0: left, 1: right)
}

// this is fundamental for this solution, a structre that is made to create a flat sorted list 
// listing all numeric values in the tree "from left to right" in the string representation by linking
// the structs in the binary tree so looping through the sorted list and manipulating the pairs is
// reflected in the binary tree and a left-to-right traversal becomes a simple loop
type vlist struct {
	pr *pair // the paire the value is stored in
	ix int   // the index of the value (0: left, 1: right)
}

// parses the string representation of a pair with nested pairs into binary tree with the root being the returned pair
// build for recursive calling thus needs to be called with up:nil / upIx:-1 for the root
func parseNum(s string, lvl int, up *pair, upIx int) (p *pair) {
	p       = &pair{}
	p.lvl   = lvl
	p.typ   = 0 
	p.up    = up
	p.upIx  = upIx
	p.prs   = make([]*pair, 2)
	p.val   = make([]int, 2)

	// detremines the position of the relevant comma and closing brackt for this pair
	cm, cl := commaClose(s)

	// uses the fact that only one digit ints are in the input to be parsed
	// NOTE: this fails if there are two digit ints in the input
	if cm == 2 {
		p.val[0]   = int(s[1] - '0')
	} else {
		p.prs[0]   = parseNum(s[1:cm], lvl+1, p, 0)
		p.typ += leftP
	}

	if cl-cm == 2 {
		p.val[1]   = int(s[cm+1] - '0')
	} else {
		p.prs[1]   = parseNum(s[cm+1:cl], lvl+1, p, 1)
		p.typ += rightP
	}
	return
}

// detects the next relevant comma (i.e. the comma that belongs to the opening bracket this string starts with)
// and detects the clsoing bracket for the opening bracket this starts with
func commaClose(s string) (cmIx, clIx int) {
	lvl := 0
	cmIx = -1
	clIx = -1

	if s[0] != '[' {
		fmt.Println("Malformed !! -", s)
	}
	for ix,c := range s {
		switch c {
		case '[':
			lvl += 1
		case ']':
			if lvl == 1 {
				clIx = ix
				return
			} 
			lvl -= 1
		case ',':
			if lvl == 1 {
				cmIx = ix
			}
		}
	}
	return
}

// debugging helper converts a binary tree into the bracket representation used in the puzzle description
func (p *pair) toString() string {
	switch p.typ {
	case noP:
		return fmt.Sprintf("[%v,%v]", p.val[0], p.val[1])
	case leftP:
		return fmt.Sprintf("[" + p.prs[0].toString() + ",%v]", p.val[1])
	case rightP:
		return fmt.Sprintf("[%v," + p.prs[1].toString() + "]", p.val[0])
	default:
		return fmt.Sprintf("[" + p.prs[0].toString() + "," + p.prs[1].toString() + "]")
	}
}

// together with the structure above, this construct makes is very easy to implement the operations
// this is a list of numeric values that is ordered by how far left values appear in the string representation
// however, the values themselves are not saved in the list but the list has pointers to the binary tree pairs
// holding the values and indicating whether the value is held in the lef tor right side of the pair
// NOTE: pairs of type 0 (two numeric values) will be linked twice in the list, once for the left and once for the right value
func (p *pair) valueList(vl []vlist) []vlist {
	switch p.typ {
	case noP: 
		return append(vl, vlist{pr:p, ix:0}, vlist{pr:p, ix:1})
	case leftP: 
		return append(p.prs[0].valueList(vl), vlist{pr:p, ix:1})
	case rightP: 
		return append([]vlist{vlist{pr:p, ix:0}}, p.prs[1].valueList(vl)...)
	case bothP:
		return append(p.prs[0].valueList(vl), p.prs[1].valueList(vl)...)
	}
	return nil
} 

// a debuggin function printing the numeric values extracted from a binary tree in order
func dump(vl []vlist) {
	for _, v := range vl {
		fmt.Printf("%v(%v) ", v.pr.val[v.ix], v.ix)
	}
	fmt.Println()
}

// once the sorted list is there, explode becomes as simple as looping through the pairs linked
// in the sorted list and once an explosion target is identified, the value right before or after
// the exploding pair is just before and after the exploding pair in this flat list. 
// returns a bool indicating whether an explosion has happened.
func (p *pair) explode() bool {

	vlst := p.valueList([]vlist{})

	for i, vl := range vlst {

		// condition for explosion
		if vl.pr.typ == 0 && vl.pr.lvl > 3 {

			if i > 0 { // not the first value?
				vlst[i-1].pr.val[vlst[i-1].ix] += vl.pr.val[0]
			}
			if (i + 2) < len(vlst) { // not the last value?
				vlst[i+2].pr.val[vlst[i+2].ix] += vl.pr.val[1]
			}

			// collapsing this pair into a '0' in the parent pair
			np := vl.pr.up
			np.val[vl.pr.upIx]  = 0
			np.prs[vl.pr.upIx]  = nil
			if vl.pr.upIx == 0 {
				np.typ -= leftP
			} else {
				np.typ -= rightP
			}

			return true
		}
	}

	return false
}

// split is also easier with the sorted list as the condition (val > 9) can be checked in the simple loop
// uses the binary tree directly to insert the new pair
func (p *pair) split() bool {

	vlst := p.valueList([]vlist{})

	for _, vl := range vlst {
		if vl.pr.val[vl.ix] > 9 {

			np := &pair{lvl:  vl.pr.lvl+1, 
			            val:  []int{vl.pr.val[vl.ix]/2, vl.pr.val[vl.ix] / 2 + (vl.pr.val[vl.ix] % 2)},
			            prs:  make([]*pair, 2),
			            typ:  0,
			            up:   vl.pr, 
			            upIx: vl.ix}

			vl.pr.prs[vl.ix] = np
			if vl.ix == 1 {
				vl.pr.typ += rightP
			} else {
				vl.pr.typ += leftP				
			}
			vl.pr.val[vl.ix] = 0
			return true
		}
	}
	return false
}

// reduction is basically trying to explode and if nothing explodes, trying to split and start again.
// If neither happend the line is stable
func (p *pair) reduce() {

	expld := true
	split := true
	for expld || split {
		expld = p.explode()
		if !expld {
			split = p.split()
		}
	} 
}

// adds one to the leve of this pair and all nested pairs
// needed as addition of two lines adds a new root pair so the existing root pairs and all nested pairs need to move down
func (p *pair) uplevel() {

	p.lvl += 1
	switch p.typ {
	case leftP: 
		p.prs[0].uplevel()
	case rightP: 
		p.prs[1].uplevel()
	case bothP:
		p.prs[0].uplevel()
		p.prs[1].uplevel()		
	}

}

// add creates a new root level pointing to the two input pairs
func (p *pair) add(p2 *pair) *pair {

	p.uplevel()
	p2.uplevel()

	np := pair{ typ:  bothP,
				lvl:  0,
				upIx: -1,
				prs:  []*pair{p, p2},
				val:  make([]int, 2)}

	p.upIx  = 0
	p.up    = &np
	p2.upIx = 1
	p2.up   = &np

	return &np
}

// simple recursive magnitude computation
func (p *pair) magnitude() int {

	switch p.typ {
	case noP:
		return p.val[0] * 3 + 2 * p.val[1]
	case leftP: 
		return p.prs[0].magnitude() * 3 + 2 * p.val[1]
	case rightP: 
		return p.val[0] * 3 + 2 * p.prs[1].magnitude()
	case bothP:
		return p.prs[0].magnitude() * 3 + 2 * p.prs[1].magnitude()
	}
	return 0
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
	input  := readTxtFile("d18." + dataset + ".txt")

	// Part 1 - read first line
	pr := parseNum(input[0], 0, nil, -1)

	// add the remaining lines one after the other
	for i:=1; i<len(input); i++ {
		pr = pr.add(parseNum(input[i], 0, nil, -1))
		pr.reduce()
	}

	// magnitude after summing all lines up
	fmt.Println("Magnitude: ", pr.magnitude())

	// Part 2 
	res := []int{} 

	// go through all combinations of lines. Since the add is not commutative
	// we allow the same pair of lines twice in either sequence
	for i1:=0; i1<len(input); i1++ {
		for i2:=0; i2<len(input); i2++ {
			p1 := parseNum(input[i1], 0, nil, -1)
			p2 := parseNum(input[i2], 0, nil, -1)
			pr = p1.add(p2)
			pr.reduce()
			res = append(res, pr.magnitude())
		}
	}
	sort.Ints(res)
	fmt.Println("Maximum magnitude: ", res[len(res)-1])

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}