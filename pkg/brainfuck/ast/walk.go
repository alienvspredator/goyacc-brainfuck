package ast

func Walk(v Visitor, node Node) {
	if v = v.Visit(node); v == nil {
		return
	}

	switch n := node.(type) {
	case *Loop:
		Walk(v, n.Body)
	case *Program:
		Walk(v, n.Body)
	case *Body:
		walkForList(v, n.List)
	}
}

func walkForList(v Visitor, list []Node) {
	for _, x := range list {
		Walk(v, x)
	}
}
