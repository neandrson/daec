package rpn

type AST struct {
	token       string
	left, right *AST
}

func NewAST(rpn *RPN) *AST {
	return new(AST)
}
