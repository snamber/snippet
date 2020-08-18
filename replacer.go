package snippet

import (
	"strings"
)

func replace(p parser, content string, s snippets) string {
	lines := strings.Split(content, "\n")

	var resultLines []string

	for _, line := range lines {
		if targetIndent, tag, ok := p.snippetLocation(line); ok {
			resultLines = append(resultLines, s[tag].String(targetIndent))
		} else {
			resultLines = append(resultLines, line)
		}
	}

	return strings.Join(resultLines, "\n")
}
