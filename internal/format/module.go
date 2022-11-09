package format

import (
	"bufio"
	"bytes"
	"fmt"
	"io"

	"github.com/busser/tftree/internal/tftree"
	"github.com/mitchellh/colorstring"
)

func Module(mod tftree.Module) string {

	c := colorstring.Colorize{
		Colors:  colorstring.DefaultColors,
		Reset:   true,
		Disable: NoColor,
	}

	var buf bytes.Buffer

	buf.WriteString(c.Color(fmt.Sprintf("[bold][blue]%s[reset] ([green]%s[reset])", mod.Name, mod.Source)))
	buf.WriteByte('\n')

	for i, child := range mod.Children {
		var childBuf bytes.Buffer
		childBuf.WriteString(Module(child))

		lastItem := i == len(mod.Children)-1
		buf.WriteString(asTree(&childBuf, lastItem))
	}

	return buf.String()
}

func asTree(r io.Reader, lastItem bool) string {
	var (
		leftRuleFirst = "├── "
		leftRule      = "│   "
	)
	if lastItem {
		leftRuleFirst = "└── "
		leftRule = "    "
	}

	var buf bytes.Buffer
	first := true

	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := sc.Text()
		if first {
			buf.WriteString(leftRuleFirst)
			first = false
		} else {
			buf.WriteString(leftRule)
		}
		buf.WriteString(line)
		buf.WriteByte('\n')
	}

	return buf.String()
}
