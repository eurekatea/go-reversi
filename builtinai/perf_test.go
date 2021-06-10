package builtinai

import (
	"math/rand"
	"testing"
)

// avg: 91.04 ns/op
func BenchmarkEvalNorm(b *testing.B) {
	bd := newBoardFromStr("+++++++++++XO++++++OOX+++OOOOXO+++OOOOOO+OXOOXXX+++OOXX++++XO+++")
	for i := 0; i < b.N; i++ {
		bd.eval(BLACK, WHITE, VALUE8x8)
	}
}

// avg: 56.4 ns/op
func BenchmarkEvalB(b *testing.B) {
	bd := newBboard8("+++++++++++XO++++++OOX+++OOOOXO+++OOOOOO+OXOOXXX+++OOXX++++XO+++")
	for i := 0; i < b.N; i++ {
		bd.eval(BLACK)
	}
}

// avg: 82 ns/op
func BenchmarkCountNorm(b *testing.B) {
	bd := newBoardFromStr("+++++++++++XO++++++OOX+++OOOOXO+++OOOOOO+OXOOXXX+++OOXX++++XO+++")
	for i := 0; i < b.N; i++ {
		_ = bd.countPieces(BLACK) - bd.countPieces(WHITE)
	}
}

// avg: 8.4 ns/op
func BenchmarkCountB(b *testing.B) {
	bd := newBboard8("+++++++++++XO++++++OOX+++OOOOXO+++OOOOOO+OXOOXXX+++OOXX++++XO+++")
	for i := 0; i < b.N; i++ {
		_ = bd.count(BLACK) - bd.count(WHITE)
	}
}

// avg: 564 ns/op 496 B/op 13 allocs/op
func BenchmarkCpy(b *testing.B) {
	bd := newBoardFromStr("+++++++++++XO++++++OOX+++OOOOXO+++OOOOOO+OXOOXXX+++OOXX++++XO+++")
	for i := 0; i < b.N; i++ {
		bd.put(WHITE, point{4, 0})
		_ = bd.Copy()
	}
}

// avg: 170 ns/op 96 B/op 2 allocs/op
func BenchmarkRevertbd(b *testing.B) {
	bd := newBoardFromStr("+++++++++XX++OOOX+++OXOO++X+XX++++++")

	for i := 0; i < b.N; i++ {
		hs := bd.put(WHITE, point{4, 0})
		bd.revert(hs)
	}
}

// avg: 40 ns/op 0 B/op 0 allocs/op
func BenchmarkCpyb(b *testing.B) {
	bd := newBboard8("+++++++++++XO++++++OOX+++OOOOXO+++OOOOOO+OXOOXXX+++OOXX++++XO+++")
	for i := 0; i < b.N; i++ {
		bd.put(WHITE, 32)
		_ = bd.cpy()
	}
}

func BenchmarkAccessNorm(b *testing.B) {
	bd := newBoardFromStr("+++++++++++XO++++++OOX+++OOOOXO+++OOOOOO+OXOOXXX+++OOXX++++XO+++")
	p := point{4, 3}
	for i := 0; i < b.N; i++ {
		_ = bd.at(p)
	}
}

// almost the same ↑↓ (0.5 ns/op)

func BenchmarkAccessB(b *testing.B) {
	bd := newBboard8("+++++++++++XO++++++OOX+++OOOOXO+++OOOOOO+OXOOXXX+++OOXX++++XO+++")
	loc := 28
	for i := 0; i < b.N; i++ {
		_ = bd.at(loc)
	}
}

// avg: 0.677 ns/op
func BenchmarkAssignNorm(b *testing.B) {
	bd := newBoardFromStr("+++++++++++XO++++++OOX+++OOOOXO+++OOOOOO+OXOOXXX+++OOXX++++XO+++")
	p := point{4, 3}
	for i := 0; i < b.N; i++ {
		bd.assign(WHITE, p)
	}
}

// avg: 1.372 ns/op
func BenchmarkAssignB(b *testing.B) {
	bd := newBboard8("+++++++++++XO++++++OOX+++OOOOXO+++OOOOOO+OXOOXXX+++OOXX++++XO+++")
	loc := 28
	for i := 0; i < b.N; i++ {
		bd.assign(WHITE, loc)
	}
}

func BenchmarkHw(b *testing.B) {
	num := rand.Uint64()
	for i := 0; i < b.N; i++ {
		hammingWeight(num)
	}
}
