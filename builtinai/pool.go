package builtinai

type pool struct {
	stack []nodes
	curr  int
}

const INIT_SIZE = 64

func newPool(cap int) pool {
	s := make([]nodes, 0)
	for i := 0; i < cap; i++ {
		s = append(s, make(nodes, 0, INIT_SIZE))
	}
	return pool{
		stack: s,
		curr:  0,
	}
}

func (p *pool) getClearOne() nodes {
	ns := p.stack[p.curr]
	p.curr++
	return ns[:0]
}

func (p *pool) freeOne() {
	p.curr--
}
