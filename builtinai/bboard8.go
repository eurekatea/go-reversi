package builtinai

const u1 uint64 = 1

var DIR = []int{-8, -7, 1, 9, 8, 7, -1, -9}

type bboard8 struct {
	black, white uint64
}

func newBboard8(input string) bboard8 {
	bd := bboard8{}
	for i := 0; i < 64; i++ {
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

func (bd bboard8) String() (res string) {
	for loc := 0; loc < 64; loc++ {
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

func (bd bboard8) visualize() (res string) {
	res = "  a b c d e f g h"
	for loc := 0; loc < 64; loc++ {
		if loc%8 == 0 {
			res += "\n" + string(rune('A'+loc/8)) + " "
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

func (bd bboard8) cpy() bboard8 {
	return bboard8{bd.black, bd.white}
}

func (bd bboard8) at(loc int) color {
	sh := u1 << loc
	if bd.black&sh != 0 {
		return BLACK
	} else if bd.white&sh != 0 {
		return WHITE
	}
	return NONE
}

func (bd *bboard8) assign(cl color, loc int) {
	sh := u1 << loc
	if cl == BLACK {
		bd.black |= sh
	} else {
		bd.white |= sh
	}
}

func (bd *bboard8) put(cl color, loc int) {
	bd.assign(cl, loc)
	bd.flip(cl, loc)
}

func (bd *bboard8) putAndCheck(cl color, loc int) bool {
	if loc < 0 {
		return false
	}
	if bd.at(loc) != NONE || !bd.isValidLoc(cl, loc) {
		return false
	}
	bd.put(cl, loc)
	bd.flip(cl, loc)
	return true
}

func (bd *bboard8) clear(loc int) {
	c := ^(u1 << loc)
	bd.black &= c
	bd.white &= c
}

func (bd *bboard8) flip(cl color, loc int) {
	var x, bounding_disk uint64
	new_disk := (u1 << loc)
	captured_disks := uint64(0)

	if cl == BLACK {
		bd.black |= new_disk

		for dir := 0; dir < 8; dir++ {
			/* Find opponent disk adjacent to the new disk. */
			x = bd.shift(new_disk, dir) & bd.white
			/* Add any adjacent opponent disk to that one, and so on. */
			x |= bd.shift(x, dir) & bd.white
			x |= bd.shift(x, dir) & bd.white
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

func (bd bboard8) allValidLoc(cl color) uint64 {
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
		x |= bd.shift(x, dir) & opp
		x |= bd.shift(x, dir) & opp

		legal |= bd.shift(x, dir) & empty
	}
	return legal
}

func (bd bboard8) hasValidMove(cl color) bool {
	return bd.allValidLoc(cl) != 0
}

func (bd bboard8) isValidLoc(cl color, loc int) bool {
	mask := u1 << loc
	return bd.allValidLoc(cl)&mask != 0
}

var (
	masks = []uint64{
		0x7F7F7F7F7F7F7F7F,
		0x007F7F7F7F7F7F7F,
		0xFFFFFFFFFFFFFFFF,
		0x00FEFEFEFEFEFEFE,
		0xFEFEFEFEFEFEFEFE,
		0xFEFEFEFEFEFEFE00,
		0xFFFFFFFFFFFFFFFF,
		0x7F7F7F7F7F7F7F00,
	}

	lshift = []uint64{
		0, 0, 0, 0, 1, 9, 8, 7,
	}

	rshift = []uint64{
		1, 9, 8, 7, 0, 0, 0, 0,
	}
)

func (bd bboard8) shift(disk uint64, dir int) uint64 {
	if dir < 8/2 {
		return (disk >> rshift[dir]) & masks[dir]
	} else {
		return (disk << lshift[dir]) & masks[dir]
	}
}

func (bd bboard8) count(cl color) int {
	var n uint64
	if cl == BLACK {
		n = bd.black
	} else {
		n = bd.white
	}
	return hammingWeight(n)
}

func (bd bboard8) emptyCount() int {
	return 64 - hammingWeight(bd.black|bd.white)
}

func (bd bboard8) isOver() bool {
	return !(bd.hasValidMove(BLACK) || bd.hasValidMove(WHITE))
}

func hammingWeight(n uint64) int {
	n = (n & 0x5555555555555555) + ((n >> 1) & 0x5555555555555555)
	n = (n & 0x3333333333333333) + ((n >> 2) & 0x3333333333333333)
	n = (n & 0x0F0F0F0F0F0F0F0F) + ((n >> 4) & 0x0F0F0F0F0F0F0F0F)
	n = (n & 0x00FF00FF00FF00FF) + ((n >> 8) & 0x00FF00FF00FF00FF)
	n = (n & 0x0000FFFF0000FFFF) + ((n >> 16) & 0x0000FFFF0000FFFF)
	n = (n & 0x00000000FFFFFFFF) + ((n >> 32) & 0x00000000FFFFFFFF)
	return int(n)
}

// var (
// 	CORNER = []uint64{
// 		0x8100000000000081,
// 		0x0042000000004200,
// 		0x0000240000240000,
// 		0x0000001818000000,
// 	}
// 	CORNERV = []int{800, -552, 62, -18}

// 	EDGE = []uint64{
// 		0x4281000000008142,
// 		0x2400810000810024,
// 		0x1800008181000018,
// 		0x0024420000422400,
// 		0x0018004242001800,
// 		0x0000182424180000,
// 	}
// 	EDGEV = []int{-286, 426, -24, -177, -82, 8}
// )

// loop unrolling
func (bd bboard8) eval(cl color) int {
	bv, wv := 0, 0
	cnt := 0

	cnt = hammingWeight(bd.black & 0x8100000000000081)
	bv += cnt * 800
	cnt = hammingWeight(bd.black & 0x0042000000004200)
	bv += cnt * -552
	cnt = hammingWeight(bd.black & 0x0000240000240000)
	bv += cnt * 62
	cnt = hammingWeight(bd.black & 0x0000001818000000)
	bv += cnt * -18
	cnt = hammingWeight(bd.black & 0x4281000000008142)
	bv += cnt * -286
	cnt = hammingWeight(bd.black & 0x2400810000810024)
	bv += cnt * 426
	cnt = hammingWeight(bd.black & 0x1800008181000018)
	bv += cnt * -24
	cnt = hammingWeight(bd.black & 0x0024420000422400)
	bv += cnt * -177
	cnt = hammingWeight(bd.black & 0x0018004242001800)
	bv += cnt * -82
	cnt = hammingWeight(bd.black & 0x0000182424180000)
	bv += cnt * 8

	cnt = hammingWeight(bd.white & 0x8100000000000081)
	wv += cnt * 800
	cnt = hammingWeight(bd.white & 0x0042000000004200)
	wv += cnt * -552
	cnt = hammingWeight(bd.white & 0x0000240000240000)
	wv += cnt * 62
	cnt = hammingWeight(bd.white & 0x0000001818000000)
	wv += cnt * -18
	cnt = hammingWeight(bd.white & 0x4281000000008142)
	wv += cnt * -286
	cnt = hammingWeight(bd.white & 0x2400810000810024)
	wv += cnt * 426
	cnt = hammingWeight(bd.white & 0x1800008181000018)
	wv += cnt * -24
	cnt = hammingWeight(bd.white & 0x0024420000422400)
	wv += cnt * -177
	cnt = hammingWeight(bd.white & 0x0018004242001800)
	wv += cnt * -82
	cnt = hammingWeight(bd.white & 0x0000182424180000)
	wv += cnt * 8

	if cl == BLACK {
		return bv - wv
	} else {
		return wv - bv
	}
}

// return the mobility (how many possible moves)
func (bd bboard8) mobility(cl color) int {
	allv := bd.allValidLoc(cl)
	return hammingWeight(allv)
}
