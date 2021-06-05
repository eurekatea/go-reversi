package builtinai

import (
	"math/rand"
	"sort"
)

type node struct {
	x, y  int
	value int
}

func newNode(x, y, value int) node {
	return node{x: x, y: y, value: value}
}

type nodes []node

func (ns nodes) Len() int {
	return len(ns)
}

func (ns nodes) Less(i, j int) bool {
	return ns[i].value > ns[j].value // descending order
}

func (ns nodes) Swap(i, j int) {
	ns[i], ns[j] = ns[j], ns[i]
}

// provide randomness
func (ns nodes) shuffle() {
	rand.Shuffle(len(ns), func(i, j int) {
		ns[i], ns[j] = ns[j], ns[i]
	})
}

func (ns nodes) sortDesc() {
	sort.Slice(ns, func(i, j int) bool {
		return ns[i].value > ns[j].value
	})
}

func (ns nodes) sortAsc() {
	sort.Slice(ns, func(i, j int) bool {
		return ns[i].value < ns[j].value
	})
}
