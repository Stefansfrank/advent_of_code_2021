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

// input parser returns the mapping algorithm alg[len = 512]
// and the image as a [][]byte with values 0 and 1
// The image is padded with two rows / columns of zeros in all directions in order to support
// the computation of all potentially changing pixels (all pixels plus one row/column on each side) without index overflows
func parseFile (lines []string) ([]byte, [][]byte) {

	// parsing the mapping algorithm from the first line
	alg := make([]byte, 512)
	for i, b := range lines[0] {
		alg[i] = byte(('.' - b) / ('.' - '#')) 
	}

	// parsing the image and padding it
	img    := make([][]byte, len(lines) + 2)
	img[0]  = make([]byte, len(lines[2]) + 4)
	img[1]  = make([]byte, len(lines[2]) + 4)
	for i := 2; i < len(lines); i++ {
		img[i] = make([]byte, len(lines[2]) + 4)
		for j, b := range lines[i] {
			img[i][j+2] = byte(('.' - b) / ('.' - '#')) 
		}
	}
	img[len(img)-2]  = make([]byte, len(lines[2]) + 4)
	img[len(img)-1]  = make([]byte, len(lines[2]) + 4)

	return alg, img
}

// calculates the new value of a bit in the image on one iteration
// the padding allows this without any index checks
func impBit(x,y int, img [][]byte, alg []byte) byte {
	
	iix := 0
	for iy := y-1; iy < y+2; iy++ {
		for ix := x-1; ix < x+2; ix++ {
			iix <<= 1 
			iix  += int(img[iy][ix])
		}
	}
	return alg[iix]
}

// prints the image and can cut the padding by 'cut' rows/columns
// only neeed for debugging
func dump(img [][]byte, cut int) {
	for y := cut; y < len(img) - cut; y++ {
		for x := cut; x < len(img[0]) - cut; x++{
			fmt.Printf("%c", '.' - img[y][x] * ('.' - '#'))
		}
		fmt.Println()
	}
	fmt.Println()
}

// counts the pixels in an image (without the padding)
func count(img [][]byte) (cnt int){

	for y := 2; y < len(img) - 2; y ++ {
		for x := 2; x < len(img[0]) - 2; x++ {
			cnt += int(img[y][x]) 
		}
	}
	return
}

// one iteration of improvement - returns the improved image in a new slice
// defOn determines whether the inifinite pixels in the new image should by 1 or 0
// (see explanation in main below as to when this is important)
func improve(img [][]byte, alg []byte, defOn bool) (nimg [][]byte) {

	// the new slice is bigger by two elements in both dimensions
	nimg    = make([][]byte, len(img) + 2)
	nimg[0] = make([]byte, len(img[0]) + 2)
	nimg[1] = make([]byte, len(img[0]) + 2)

	// looping through the input image including all values in the first row/column of padding
	for y := 1; y < len(img) - 1; y++ {
		nimg[y+1] = make([]byte, len(img[0]) + 2)
		for x := 1; x < len(img[0]) - 1; x++ {
			nimg[y+1][x+1] = impBit(x, y, img, alg)
		}
	}

	nimg[len(nimg) - 2] = make([]byte, len(img[0]) + 2)
	nimg[len(nimg) - 1] = make([]byte, len(img[0]) + 2)

	// fill all the padding with 1 if requested
	if defOn {
		for i := 0; i < len(nimg[0]); i++ {
			nimg[0][i] = 1
			nimg[1][i] = 1
			nimg[len(nimg)-2][i] = 1
			nimg[len(nimg)-1][i] = 1
		}

		for i := 2; i < len(nimg); i ++ {
			nimg[i][0] = 1
			nimg[i][1] = 1
			nimg[i][len(nimg[0]) - 2] = 1
			nimg[i][len(nimg[0]) - 1] = 1
		}
	}
	return nimg
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

	start    := time.Now()
	input    := readTxtFile("d20." + dataset + ".txt")
	alg, img := parseFile(input)

	// detect the pattern for the value of the infinite pixels:
	// - they are always 0 if alg maps 000000000 to 0 (Pattern 0)
	// - they are alternating 0 and 1 if alg maps 000000000 to 1 and 111111111 to 0 (Pattern 1)
	// - they are always 1 after the first iteration if alg maps both 000000000 and 111111111 to 1 (Pattern 2)
	infPat := 0
	if alg[0] == 1 {
		if alg[len(alg) - 1] == 1 {
			infPat = 2
		} else {
			infPat = 1
		}
	} 

	for i := 0; i < 50; i++ {
		img = improve(img, alg, (i%2 == 0) && (infPat == 1) || infPat == 2)	
		if i == 1 || i == 49 {
			fmt.Println("Pattern after", i+1, "iterations:", count(img))
		}
	}

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}