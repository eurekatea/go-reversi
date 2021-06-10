package builtinai

type bboard6 struct {
	black, white uint64
}

func newBboard6(input string) bboard6 {
	bd := bboard6{}
	for i := 0; i < 36; i++ {
		switch input[i] {
		case 'X':
			bd.assign(BLACK, i)
		case 'O':
			bd.assign(WHITE, i)
		case '+':
		default:
			panic("input err: " + string(input[i]))
		}
	}
	return bd
}

func (bd bboard6) String() (res string) {
	for loc := 0; loc < 36; loc++ {
		switch bd.at(loc) {
		case NONE:
			res += "+"
		case BLACK:
			res += "X"
		case WHITE:
			res += "O"
		default:
			panic("err: " + bd.at(loc).String())
		}
	}
	return
}

func (bd bboard6) visualize() (res string) {
	res = "  a b c d e f"
	for loc := 0; loc < 36; loc++ {
		if loc%6 == 0 {
			res += "\n" + string(rune('A'+loc/6)) + " "
		}
		switch bd.at(loc) {
		case NONE:
			res += "+ "
		case BLACK:
			res += "X "
		case WHITE:
			res += "O "
		default:
			panic("err: " + bd.at(loc).String())
		}
	}
	return res + "\n"
}

func (bd bboard6) cpy() bboard6 {
	return bboard6{bd.black, bd.white}
}

func (bd bboard6) at(loc int) color {
	sh := u1 << loc
	if bd.black&sh != 0 {
		return BLACK
	} else if bd.white&sh != 0 {
		return WHITE
	}
	return NONE
}

func (bd *bboard6) assign(cl color, loc int) {
	sh := u1 << loc
	if cl == BLACK {
		bd.black |= sh
	} else {
		bd.white |= sh
	}
}

func (bd *bboard6) put(cl color, loc int) {
	bd.assign(cl, loc)
	bd.flip(cl, loc)
}

func (bd *bboard6) putAndCheck(cl color, loc int) bool {
	if loc < 0 || loc >= 36 {
		return false
	}
	if bd.at(loc) != NONE || !bd.isValidLoc(cl, loc) {
		return false
	}
	bd.put(cl, loc)
	bd.flip(cl, loc)
	return true
}

func (bd *bboard6) clear(loc int) {
	c := ^(u1 << loc)
	bd.black &= c
	bd.white &= c
}

func (bd *bboard6) flip(cl color, loc int) {
	var x, bounding_disk uint64
	new_disk := (u1 << loc)
	captured_disks := uint64(0)

	if cl == BLACK {
		bd.black |= new_disk

		for dir := 0; dir < 8; dir++ {
			// find opponent disk adjacent to the new disk.
			x = bd.shift(new_disk, dir) & bd.white
			/* Add any adjacent opponent disk to that one, and so on. */
			x |= bd.shift(x, dir) & bd.white
			x |= bd.shift(x, dir) & bd.white
			x |= bd.shift(x, dir) & bd.white
			/* Determine whether the disks were captured. */
			bounding_disk = bd.shift(x, dir) & bd.black

			if bounding_disk != 0 {
				captured_disks |= x
			}
		}
		bd.black ^= captured_disks
		bd.white ^= captured_disks
	} else {
		bd.white |= new_disk

		for dir := 0; dir < 8; dir++ {
			/* Find opponent disk adjacent to the new disk. */
			x = bd.shift(new_disk, dir) & bd.black
			/* Add any adjacent opponent disk to that one, and so on. */
			x |= bd.shift(x, dir) & bd.black
			x |= bd.shift(x, dir) & bd.black
			x |= bd.shift(x, dir) & bd.black
			/* Determine whether the disks were captured. */
			bounding_disk = bd.shift(x, dir) & bd.white

			if bounding_disk != 0 {
				captured_disks |= x
			}
		}
		bd.white ^= captured_disks
		bd.black ^= captured_disks
	}
}

func (bd bboard6) allValidLoc(cl color) uint64 {
	var legal uint64
	var self, opp uint64

	if cl == BLACK {
		self = bd.black
		opp = bd.white
	} else {
		self = bd.white
		opp = bd.black
	}
	empty := ^(self | opp)

	for dir := 0; dir < 8; dir++ {
		x := bd.shift(self, dir) & opp
		x |= bd.shift(x, dir) & opp
		x |= bd.shift(x, dir) & opp
		x |= bd.shift(x, dir) & opp

		legal |= bd.shift(x, dir) & empty
	}
	return legal
}

func (bd bboard6) hasValidMove(cl color) bool {
	return bd.allValidLoc(cl) != 0
}

func (bd bboard6) isValidLoc(cl color, loc int) bool {
	mask := u1 << loc
	return bd.allValidLoc(cl)&mask != 0
}

// var (
// 	masks6 = []uint64{
// 		0x7DF7DF7DF, // e
// 		0x01F7DF7DF, // se
// 		0xFFFFFFFFF, // s
// 		0x03EFBEFBE, // sw
// 		0xFBEFBEFBE, // w
// 		0xFBEFBEF80, // nw
// 		0xFFFFFFFFF, // n
// 		0x7DF7DF7C0,
// 	}

// 	lshift6 = []uint64{
// 		0, 0, 0, 0, 1, 7, 6, 5,
// 	}

// 	rshift6 = []uint64{
// 		1, 7, 6, 5, 0, 0, 0, 0,
// 	}
// )

func (bd bboard6) shift(disk uint64, dir int) uint64 {
	switch dir {
	case 0:
		return (disk >> 1) & 0x7DF7DF7DF // e
	case 1:
		return (disk >> 7) & 0x01F7DF7DF // se
	case 2:
		return (disk >> 6) & 0xFFFFFFFFF // s
	case 3:
		return (disk >> 5) & 0x03EFBEFBE // sw
	case 4:
		return (disk << 1) & 0xFBEFBEFBE // w
	case 5:
		return (disk << 7) & 0xFBEFBEF80 // nw
	case 6:
		return (disk << 6) & 0xFFFFFFFFF // n
	case 7:
		return (disk << 5) & 0x7DF7DF7C0 // ne
	}
	panic("dir error")
}

func (bd bboard6) count(cl color) int {
	if cl == BLACK {
		return hammingWeight(bd.black)
	} else {
		return hammingWeight(bd.white)
	}
}

func (bd bboard6) emptyCount() int {
	return 36 - hammingWeight(bd.black|bd.white)
}

func (bd bboard6) isOver() bool {
	return !(bd.hasValidMove(BLACK) || bd.hasValidMove(WHITE))
}

// loop unrolling
func (bd bboard6) eval(cl color) int {
	bv, wv := 0, 0
	cnt := 0

	cnt = hammingWeight(bd.black & 0x840000021)
	bv += cnt * 100
	cnt = hammingWeight(bd.black & 0x4A1000852)
	bv += cnt * -36
	cnt = hammingWeight(bd.black & 0x30086100C)
	bv += cnt * 53
	cnt = hammingWeight(bd.black & 0x012000480)
	bv += cnt * -69
	cnt = hammingWeight(bd.black & 0x00C492300)
	bv += cnt * -10
	cnt = hammingWeight(bd.black & 0x00030C000)
	bv += cnt * -2

	cnt = hammingWeight(bd.white & 0x840000021)
	wv += cnt * 100
	cnt = hammingWeight(bd.white & 0x4A1000852)
	wv += cnt * -36
	cnt = hammingWeight(bd.white & 0x30086100C)
	wv += cnt * 53
	cnt = hammingWeight(bd.white & 0x012000480)
	wv += cnt * -69
	cnt = hammingWeight(bd.white & 0x00C492300)
	wv += cnt * -10
	cnt = hammingWeight(bd.white & 0x00030C000)
	wv += cnt * -2

	if cl == BLACK {
		return bv - wv
	} else {
		return wv - bv
	}
}

// return the mobility (how many possible moves)
func (bd bboard6) mobility(cl color) int {
	allv := bd.allValidLoc(cl)
	return hammingWeight(allv)
}
