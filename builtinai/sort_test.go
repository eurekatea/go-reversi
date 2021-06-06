package builtinai

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
	"time"
)

func qsortSimple(s []node) {
	if len(s) < 2 {
		return
	}

	left, right := 0, len(s)-1

	pivot := 0

	s[pivot], s[right] = s[right], s[pivot]

	for i := range s {
		if s[i].value < s[right].value {
			s[left], s[i] = s[i], s[left]
			left++
		}
	}

	s[left], s[right] = s[right], s[left]

	qsortSimple(s[:left])
	qsortSimple(s[left+1:])
}

// avg: 180 ns/op 96 B/op 3 allocs/op
func BenchmarkSortBuiltIn(b *testing.B) {
	ns := make(nodes, 10)
	rand.Seed(time.Now().Unix())
	for i := range ns {
		ns[i] = newNode(rand.Int(), rand.Int(), rand.Intn(10))
	}

	for i := 0; i < b.N; i++ {
		sort.Slice(ns, func(i, j int) bool {
			return ns[i].value < ns[j].value
		})
	}
	for i := 0; i < 10; i++ {
		fmt.Print(ns[i].value, " ")
	}
	fmt.Println()
}

// 2 times faster than standard libary sort (on small slices)
// much much slower on large slices
// avg: 90 ns/op 0 B/op 0 allocs/op
func BenchmarkQSortSimple(b *testing.B) {
	n := make(nodes, 10)
	rand.Seed(time.Now().Unix())
	for i := range n {
		n[i] = newNode(rand.Int(), rand.Int(), rand.Intn(10))
	}

	for i := 0; i < b.N; i++ {
		qsortSimple(n)
	}
	for i := 0; i < 10; i++ {
		fmt.Print(n[i].value, " ")
	}
	fmt.Println()
}

// modified from standard libary
// avg: 16.5 ns/op 0 B/op 0 allocs/op (10 times faster before modified)
// fastest
func BenchmarkQSortModified(b *testing.B) {
	n := make(nodes, 10)
	rand.Seed(time.Now().Unix())
	for i := range n {
		n[i] = newNode(rand.Int(), rand.Int(), rand.Intn(10))
	}

	for i := 0; i < b.N; i++ {
		qsortDesc(n)
	}
	for i := 0; i < 10; i++ {
		fmt.Print(n[i].value, " ")
	}
	fmt.Println()
}
