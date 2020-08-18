package snippet

import (
	"fmt"
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

type metadata struct {
	sourceIndent string
	tag          string
}

type snippet struct {
	meta    metadata
	content []string
}

func (s *snippet) addLine(line string) {
	s.content = append(s.content, line)
}

func (s snippet) String(targetIndent string) string {
	var result []string
	for _, line := range s.content {
		result = append(result, targetIndent+strings.TrimPrefix(line, s.meta.sourceIndent))
	}
	return strings.Join(result, "\n")
}

type snippets map[string]snippet

var (
	errEndNotSpecified      = errors.New("END SNIPPET undefined")
	errNoParserForExtension = errors.New("no comment prefix defined for file extension")
)

func snippetsFromContent(p parser, content string) (snippets, error) {
	ss := snippets{}
	inSnippet := false
	activeSnippet := snippet{}
	for _, line := range strings.Split(content, "\n") {
		if p.snippetEnd(line) {
			inSnippet = false
			ss.addSnippet(activeSnippet)
			activeSnippet = snippet{}
		}
		if inSnippet {
			activeSnippet.addLine(line)
		}
		if m, ok := p.snippetStart(line); ok {
			inSnippet = true
			activeSnippet.meta = m
		}
	}
	if inSnippet {
		return nil, errors.Wrapf(errEndNotSpecified,
			"can not extract snippets from file, end with '%s END'", p.commentPrefix)
	}
	return ss, nil
}

func (s snippets) addSnippet(newSnippets ...snippet) {
	for _, snp := range newSnippets {
		s[snp.meta.tag] = snp
	}
}

func (s snippets) addSnippets(newSnippets snippets) {
	for key, value := range newSnippets {
		s[key] = value
	}
}

type parser struct {
	start         *regexp.Regexp
	end           *regexp.Regexp
	location      *regexp.Regexp
	commentPrefix string
}

// snippetStart parses a line, and returns the snippet metadata
// and an "ok value" if the line was the start of a snippet or not
func (p parser) snippetStart(line string) (metadata, bool) {
	captures := p.start.FindStringSubmatch(line)
	if len(captures) == 3 {
		return metadata{
			sourceIndent: captures[1],
			tag:          captures[2],
		}, true
	}
	return metadata{}, false
}

// snippetEnd parses a line, and returns an "ok value" if the line was the end of a snippet or not
func (p parser) snippetEnd(line string) bool {
	found := p.end.FindString(line)
	return len(found) > 0
}

// snippetLocation parses a line, and returns the indentation, the snippet tag
// and true if the snippet target location is that line
func (p parser) snippetLocation(line string) (string, string, bool) {
	captures := p.location.FindStringSubmatch(line)
	if len(captures) == 3 {
		return captures[1], captures[2], true
	}
	return "", "", false
}

func newParser(extension string) (parser, error) {
	prefix, ok := commentPrefixMap[extension]
	if !ok {
		return parser{}, errors.Wrapf(errNoParserForExtension, "cannot create parser for file extension %s", extension)
	}

	startStr := fmt.Sprintf(`^([ \t]*)%s START SNIPPET ([a-zA-Z0-9_-]+)[ \t]*$`, prefix)
	endStr := fmt.Sprintf(`^[ \t]*%s END SNIPPET[ \t]*$`, prefix)
	locationStr := fmt.Sprintf(`^([ \t]*)%s PUT SNIPPET ([a-zA-Z0-9_-]+)[ \t]*$`, prefix)

	return parser{
		start:         regexp.MustCompile(startStr),
		end:           regexp.MustCompile(endStr),
		location:      regexp.MustCompile(locationStr),
		commentPrefix: prefix}, nil
}
