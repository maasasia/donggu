package code

import (
	"fmt"
	"io"
)

type icbIndentCommand struct{}
type icbUnindentCommand struct{}

type IndentedCodeBuilder struct {
	Commands []interface{}
}

func (i *IndentedCodeBuilder) Indent() {
	i.Commands = append(i.Commands, icbIndentCommand{})
}

func (i *IndentedCodeBuilder) Unindent() {
	i.Commands = append(i.Commands, icbUnindentCommand{})
}

func (i *IndentedCodeBuilder) AppendLines(lines ...string) {
	for _, line := range lines {
		i.Commands = append(i.Commands, line)
	}
}

func (i *IndentedCodeBuilder) IndentedLines(lines ...string) {
	i.Commands = append(i.Commands, icbIndentCommand{})
	for _, line := range lines {
		i.Commands = append(i.Commands, line)
	}
	i.Commands = append(i.Commands, icbUnindentCommand{})
}

func (i *IndentedCodeBuilder) AppendBlock(block IndentedCodeBuilder) {
	i.Commands = append(i.Commands, block.Commands...)
}

func (i *IndentedCodeBuilder) IndentedBlock(block IndentedCodeBuilder) {
	i.Commands = append(i.Commands, icbIndentCommand{})
	i.Commands = append(i.Commands, block.Commands...)
	i.Commands = append(i.Commands, icbUnindentCommand{})
}

func (i *IndentedCodeBuilder) Build(w io.Writer) {
	indentLevel := ""
	for _, cmd := range i.Commands {
		if _, ok := cmd.(icbIndentCommand); ok {
			indentLevel += "  "
		} else if _, ok := cmd.(icbUnindentCommand); ok {
			indentLevel = indentLevel[0 : len(indentLevel)-2]
		} else {
			w.Write([]byte(fmt.Sprintf("%s%s\n", indentLevel, cmd)))
		}
	}
}
