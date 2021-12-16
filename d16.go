package main

import (
	"fmt"
	"time"
	"encoding/hex"
)

// returns the bit at position 'pos' in a byte array
func bit(msg []byte, pos int) byte {
	byix := pos / 8
	biix := 7 - pos % 8
	return (msg[byix] & (1 << biix)) >> biix
}

// returns the number starting at bit position 'from' and is 'cnt' bits long
func bits(msg []byte, from, cnt int) (num int) {
	for i := 0; i < cnt; i++ {
		num <<=1
		num += int(bit(msg, from+i))
	}
	return
}

// the core packet structure
type packet struct {
	vrs  int // version
	typ  int // type of packet	
	len  int // the length of the packet in bits
	val  int // the computed value of the packet
	vSum int // version numbers sum of this packet and all sub-packets
	subs []packet // the list of sub-packets
}

// parses a package from byte array 'msg' starting at bit position 'from'
// this is designed to be used in a recursive way
func parse(msg []byte, from int) (pkt packet) {

	ix      := from // index of the parsing as it goes through
	pkt.vrs  = bits(msg, ix, 3)
	pkt.typ  = bits(msg, ix + 3, 3)
	ix += 6
	pkt.vSum = pkt.vrs
	pkt.subs = []packet{}

	switch pkt.typ {
	
	// value / constant
	case 4:
		cont := 1
		num  := 0
		for (cont == 1) {
			num <<= 4
			cont  = bits(msg, ix, 1)
			num  += bits(msg, ix + 1, 4)
			ix   += 5
		}
		pkt.val  = num
		pkt.len  = ix - from
		return
	
	// operation
	default:
		tln  := bits(msg, ix, 1) // type of parameter limit
		ix += 1

		if tln == 0 {
			subLmt := bits(msg, ix, 15) + ix + 15
			ix     += 15
			for (ix < subLmt) {
				npkt     := parse(msg, ix)
				pkt.subs  = append(pkt.subs, npkt)
				pkt.vSum += npkt.vSum
				ix       += npkt.len
			}
			pkt.len = ix - from

		} else {
			subIx := bits(msg, ix, 11)
			ix  += 11
			for i := 0; i < subIx; i++ {
				npkt     := parse(msg, ix)
				pkt.subs  = append(pkt.subs, npkt)
				pkt.vSum += npkt.vSum
				ix     += npkt.len
			}
			pkt.len = ix - from

		}

		// with all subs recursively resolved and added, evaluate operation
		pkt.val = pkt.eval()
	}
	return
}

// evaluates the operation on a packet
func (p packet) eval() (val int) {
	switch p.typ {

	// sum
	case 0:
		for _, pkt := range p.subs {
			val += pkt.val
		}
	
	// multiplication
	case 1:
		val = 1
		for _, pkt := range p.subs {
			val *= pkt.val
		}

	// min
	case 2:
		val = p.subs[0].val 
		for i := 1; i < len(p.subs); i++ {
			if p.subs[i].val < val {
				val = p.subs[i].val
			}
		}

	// max
	case 3:
		val = p.subs[0].val 
		for i := 1; i < len(p.subs); i++ {
			if p.subs[i].val > val {
				val = p.subs[i].val
			}
		}

	// greater than
	case 5:
		if p.subs[0].val > p.subs[1].val {
			val = 1 
		} else {
			val = 0
		}

	// less than
	case 6:
		if p.subs[0].val < p.subs[1].val {
			val = 1 
		} else {
			val = 0
		}

	// equals
	case 7:
		if p.subs[0].val == p.subs[1].val {
			val = 1 
		} else {
			val = 0
		}
	}

	return
}

// print out packet for debugging
// 'lvl' controls the intendation level (as this is recursive)
func (p packet) dump(lvl int) {

	for i := 0; i < lvl; i++ {
		fmt.Print("\t")		
	}
	fmt.Printf("Packet - Type %v, v%v, Value %v, Bitlen %v, vSum %v\n", p.typ, p.vrs, p.val, p.len, p.vSum)
	for _,pkt := range p.subs {
		pkt.dump(lvl+1)
	}
}

// input data and examples
func input() []string {
	inp := make([]string, 16)
	inp[0] = "620D49005AD2245800D0C9E72BD279CAFB0016B1FA2B1802DC00D0CC611A47FCE2A4ACE1DD144BFABBFACA002FB2C6F33DFF4A0C0119B169B013005F003720004263644384800087C3B8B51C26B449130802D1A0068A5BD7D49DE793A48B5400D8293B1F95C5A3005257B880F5802A00084C788AD0440010F8490F608CACE034401AB4D0F5802726B3392EE2199628CEA007001884005C92015CC8051800130EC0468A01042803B8300D8E200788018C027890088CE0049006028012AB00342A0060801B2EBE400424933980453EFB2ABB36032274C026E4976001237D964FF736AFB56F254CB84CDF136C1007E7EB42298FE713749F973F7283005656F902A004067CD27CC1C00D9CB5FDD4D0014348010C8331C21710021304638C513006E234308B060094BEB76CE3966AA007C6588A5670DC3754395485007A718A7F149CA2DD3B6E7B777800118E7B59C0ECF5AE5D3B6CB1496BAE53B7ADD78C013C00CD2629BF5371D1D4C537EA6E3A3E95A3E180592AC7246B34032CF92804001A1CCF9BA521782ECBD69A98648BC18025800F8C9C37C827CA7BEFB31EADF0AE801BA42B87935B8EF976194EEC426AAF640168CECAF84BC004AE7D1673A6A600B4AB65802D230D35CF81B803D3775683F3A3860087802132FB32F322C92A4C402524F2DE006E8000854378F710C0010D8F30FE224AE428C015E00D40401987F06E3600021D0CE3EC228DA000574E4C3080182931E936E953B200BF656E15400D3496E4A725B92998027C00A84EEEE6B347D30BE60094E537AA73A1D600B880371AA36C3200043235C4C866C018E4963B7E7AA2B379918C639F1550086064BB148BA499EC731004E1AC966BDBC7646600C080370822AC4C1007E38C428BE0008741689D0ECC01197CF216EA16802D3748FE91B25CAF6D5F11C463004E4FD08FAF381F6004D3232CC93E7715B463F780"

	// part 1 examples
	inp[1] = "D2FE28"
	inp[2] = "38006F45291200"
	inp[3] = "EE00D40C823060"
	inp[4] = "8A004A801A8002F478"
	inp[5] = "620080001611562C8802118E34"
	inp[6] = "C0015000016115A2E0802F182340"
	inp[7] = "A0016C880162017C3686B18A3D4780"

	// part 2 examples
	inp[8] = "C200B40A82"
	inp[9] = "04005AC33890"
	inp[10] = "880086C3E88112"
	inp[11] = "CE00C43D881120"
	inp[12] = "D8005AC2A8F0"
	inp[13] = "F600BC2D8F"
	inp[14] = "9C005AC2F8F0"
	inp[15] = "9C0141080250320F1802104A08"	
	return inp
}

// ------- MAIN -------
func main() {

	start := time.Now()
	inp := input()

	msg, _ := hex.DecodeString(inp[0])
	pp := parse(msg, 0)
	fmt.Println("Version sum:",pp.vSum)
	fmt.Println("Value:",pp.val)

	fmt.Println("Execution time:", time.Since(start))
}