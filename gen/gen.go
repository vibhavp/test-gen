package gen

import (
	"io"
	"strings"
)

type Var struct {
	name string
}

type GenContext struct {
	w io.Writer
	currentIndentLevel int
	variables []Var
}

func NewGenContext(w io.Writer) *GenContext {
	return &GenContext{
		w: w,
		currentIndentLevel: 0,
	}
}

func (g *GenContext) WriteLine(s string) {
	var indent string
	if g.currentIndentLevel != 0 {
		indent = strings.Repeat(" ", 4 * g.currentIndentLevel)
	}

	g.w.Write([]byte(indent + s + "\n"))
}
