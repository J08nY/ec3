package gocode

import (
	"fmt"
	"go/format"
	"go/types"

	"github.com/mmcloughlin/ec3/internal/print"
)

type Generator struct {
	print.Buffer
}

func NewGenerator() Generator {
	return Generator{
		Buffer: print.NewBuffer(),
	}
}

func (g *Generator) Package(name string) {
	g.Linef("package %s", name)
}

// Comment writes comment lines prefixed with "// ".
func (g *Generator) Comment(lines ...string) {
	for _, line := range lines {
		g.Linef("// %s", line)
	}
}

func (g *Generator) Commentf(format string, args ...interface{}) {
	g.Comment(fmt.Sprintf(format, args...))
}

func (g *Generator) CodeGenerationWarning(by string) {
	g.Commentf("Code generated by %s. DO NOT EDIT.", by)
	g.NL() // newline to ensure it does not get attached to package documentation
}

func (g *Generator) EnterBlock() {
	g.Linef("{")
	g.Indent()
}

func (g *Generator) LeaveBlock() {
	g.Dedent()
	g.Linef("}")
}

func (g *Generator) Function(name string, s *types.Signature) {
	g.Printf("func %s", name)
	types.WriteSignature(g.Buf, s, nil)
	g.EnterBlock()
}

func (g *Generator) Formatted() ([]byte, error) {
	b, err := g.Result()
	if err != nil {
		return nil, err
	}
	return format.Source(b)
}
