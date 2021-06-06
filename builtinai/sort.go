// modified from standard libary "sort"
// Copyright 2017 The Go Authors. All rights reserved.

package builtinai

func maxDepth(n int) int {
	var depth int
	for i := n; i > 0; i >>= 1 {
		depth++
	}
	return depth * 2
}

func qsortAsc(data nodes) {
	length := len(data)
	quickSort_func_Asc(data, 0, length, maxDepth(length))
}

func quickSort_func_Asc(data nodes, a, b, maxDepth int) {
	for b-a > 12 {
		if maxDepth == 0 {
			heapSort_func_Asc(data, a, b)
			return
		}
		maxDepth--
		mlo, mhi := doPivot_func_Asc(data, a, b)
		if mlo-a < b-mhi {
			quickSort_func_Asc(data, a, mlo, maxDepth)
			a = mhi
		} else {
			quickSort_func_Asc(data, mhi, b, maxDepth)
			b = mlo
		}
	}
	if b-a > 1 {
		for i := a + 6; i < b; i++ {
			if data.Less(i, i-6) {
				data.Swap(i, i-6)
			}
		}
		insertionSort_func_Asc(data, a, b)
	}
}

func insertionSort_func_Asc(data nodes, a, b int) {
	for i := a + 1; i < b; i++ {
		for j := i; j > a && data.Less(j, j-1); j-- {
			data.Swap(j, j-1)
		}
	}
}

func medianOfThree_func_Asc(data nodes, m1, m0, m2 int) {
	if data.Less(m1, m0) {
		data.Swap(m1, m0)
	}
	if data.Less(m2, m1) {
		data.Swap(m2, m1)
		if data.Less(m1, m0) {
			data.Swap(m1, m0)
		}
	}
}

func doPivot_func_Asc(data nodes, lo, hi int) (midlo, midhi int) {
	m := int(uint(lo+hi) >> 1)
	if hi-lo > 40 {
		s := (hi - lo) / 8
		medianOfThree_func_Asc(data, lo, lo+s, lo+2*s)
		medianOfThree_func_Asc(data, m, m-s, m+s)
		medianOfThree_func_Asc(data, hi-1, hi-1-s, hi-1-2*s)
	}
	medianOfThree_func_Asc(data, lo, m, hi-1)
	pivot := lo
	a, c := lo+1, hi-1
	for ; a < c && data.Less(a, pivot); a++ {
	}
	b := a
	for {
		for ; b < c && !data.Less(pivot, b); b++ {
		}
		for ; b < c && data.Less(pivot, c-1); c-- {
		}
		if b >= c {
			break
		}
		data.Swap(b, c-1)
		b++
		c--
	}
	protect := hi-c < 5
	if !protect && hi-c < (hi-lo)/4 {
		dups := 0
		if !data.Less(pivot, hi-1) {
			data.Swap(c, hi-1)
			c++
			dups++
		}
		if !data.Less(b-1, pivot) {
			b--
			dups++
		}
		if !data.Less(m, pivot) {
			data.Swap(m, b-1)
			b--
			dups++
		}
		protect = dups > 1
	}
	if protect {
		for {
			for ; a < b && !data.Less(b-1, pivot); b-- {
			}
			for ; a < b && data.Less(a, pivot); a++ {
			}
			if a >= b {
				break
			}
			data.Swap(a, b-1)
			a++
			b--
		}
	}
	data.Swap(pivot, b-1)
	return b - 1, c
}

func siftDown_func_Asc(data nodes, lo, hi, first int) {
	root := lo
	for {
		child := 2*root + 1
		if child >= hi {
			break
		}
		if child+1 < hi && data.Less(first+child, first+child+1) {
			child++
		}
		if !data.Less(first+root, first+child) {
			return
		}
		data.Swap(first+root, first+child)
		root = child
	}
}

func heapSort_func_Asc(data nodes, a, b int) {
	first := a
	lo := 0
	hi := b - a
	for i := (hi - 1) / 2; i >= 0; i-- {
		siftDown_func_Asc(data, i, hi, first)
	}
	for i := hi - 1; i >= 0; i-- {
		data.Swap(first, first+i)
		siftDown_func_Asc(data, lo, i, first)
	}
}

func qsortDesc(data nodes) {
	length := len(data)
	quickSort_func_Desc(data, 0, length, maxDepth(length))
}

func quickSort_func_Desc(data nodes, a, b, maxDepth int) {
	for b-a > 12 {
		if maxDepth == 0 {
			heapSort_func_Desc(data, a, b)
			return
		}
		maxDepth--
		mlo, mhi := doPivot_func_Desc(data, a, b)
		if mlo-a < b-mhi {
			quickSort_func_Desc(data, a, mlo, maxDepth)
			a = mhi
		} else {
			quickSort_func_Desc(data, mhi, b, maxDepth)
			b = mlo
		}
	}
	if b-a > 1 {
		for i := a + 6; i < b; i++ {
			if data.Large(i, i-6) {
				data.Swap(i, i-6)
			}
		}
		insertionSort_func_Desc(data, a, b)
	}
}

func insertionSort_func_Desc(data nodes, a, b int) {
	for i := a + 1; i < b; i++ {
		for j := i; j > a && data.Large(j, j-1); j-- {
			data.Swap(j, j-1)
		}
	}
}

func medianOfThree_func_Desc(data nodes, m1, m0, m2 int) {
	if data.Large(m1, m0) {
		data.Swap(m1, m0)
	}
	if data.Large(m2, m1) {
		data.Swap(m2, m1)
		if data.Large(m1, m0) {
			data.Swap(m1, m0)
		}
	}
}

func doPivot_func_Desc(data nodes, lo, hi int) (midlo, midhi int) {
	m := int(uint(lo+hi) >> 1)
	if hi-lo > 40 {
		s := (hi - lo) / 8
		medianOfThree_func_Desc(data, lo, lo+s, lo+2*s)
		medianOfThree_func_Desc(data, m, m-s, m+s)
		medianOfThree_func_Desc(data, hi-1, hi-1-s, hi-1-2*s)
	}
	medianOfThree_func_Desc(data, lo, m, hi-1)
	pivot := lo
	a, c := lo+1, hi-1
	for ; a < c && data.Large(a, pivot); a++ {
	}
	b := a
	for {
		for ; b < c && !data.Large(pivot, b); b++ {
		}
		for ; b < c && data.Large(pivot, c-1); c-- {
		}
		if b >= c {
			break
		}
		data.Swap(b, c-1)
		b++
		c--
	}
	protect := hi-c < 5
	if !protect && hi-c < (hi-lo)/4 {
		dups := 0
		if !data.Large(pivot, hi-1) {
			data.Swap(c, hi-1)
			c++
			dups++
		}
		if !data.Large(b-1, pivot) {
			b--
			dups++
		}
		if !data.Large(m, pivot) {
			data.Swap(m, b-1)
			b--
			dups++
		}
		protect = dups > 1
	}
	if protect {
		for {
			for ; a < b && !data.Large(b-1, pivot); b-- {
			}
			for ; a < b && data.Large(a, pivot); a++ {
			}
			if a >= b {
				break
			}
			data.Swap(a, b-1)
			a++
			b--
		}
	}
	data.Swap(pivot, b-1)
	return b - 1, c
}

func siftDown_func_Desc(data nodes, lo, hi, first int) {
	root := lo
	for {
		child := 2*root + 1
		if child >= hi {
			break
		}
		if child+1 < hi && data.Large(first+child, first+child+1) {
			child++
		}
		if !data.Large(first+root, first+child) {
			return
		}
		data.Swap(first+root, first+child)
		root = child
	}
}

func heapSort_func_Desc(data nodes, a, b int) {
	first := a
	lo := 0
	hi := b - a
	for i := (hi - 1) / 2; i >= 0; i-- {
		siftDown_func_Desc(data, i, hi, first)
	}
	for i := hi - 1; i >= 0; i-- {
		data.Swap(first, first+i)
		siftDown_func_Desc(data, lo, i, first)
	}
}
