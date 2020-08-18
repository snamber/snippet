package snippet

import (
	"reflect"
	"regexp"
	"testing"
)

func Test_snippets_Add(t *testing.T) {
	s := snippets{
		"foo": snippet{
			meta:    metadata{tag: "foo"},
			content: []string{"some content"},
		},
	}
	newSnippets := []snippet{{
		content: []string{"\t\tsome more content"},
		meta: metadata{
			sourceIndent: "\t",
			tag:          "bar",
		},
	}, {
		content: []string{"  yet more content"},
		meta: metadata{
			sourceIndent: "  ",
			tag:          "baz",
		},
	}}

	s.addSnippet(newSnippets...)

	want := snippets{
		"foo": snippet{
			content: []string{"some content"},
			meta:    metadata{tag: "foo"}},
		"bar": snippet{
			content: []string{"\t\tsome more content"},
			meta: metadata{
				sourceIndent: "\t",
				tag:          "bar",
			}},
		"baz": snippet{
			content: []string{"  yet more content"},
			meta: metadata{
				sourceIndent: "  ",
				tag:          "baz",
			}},
	}

	if !reflect.DeepEqual(s, want) {
		t.Errorf("got %v, want %v", s, want)
	}
}

func Test_parser_snippetStart(t *testing.T) {
	goParser, _ := newParser(".go")

	type fields struct {
		Regexp *regexp.Regexp
	}
	type args struct {
		line string
	}
	tests := []struct {
		name   string
		parser parser
		args   args
		want   metadata
		want1  bool
	}{
		{
			name:   "positive test",
			parser: goParser,
			args:   args{line: "  // START SNIPPET foo-bar_baz123"},
			want: metadata{
				sourceIndent: "  ",
				tag:          "foo-bar_baz123",
			},
			want1: true,
		}, {
			name:   "white-space at the end",
			parser: goParser,
			args:   args{line: "  // START SNIPPET foo-bar_baz123  \t"},
			want: metadata{
				sourceIndent: "  ",
				tag:          "foo-bar_baz123",
			},
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.parser
			got, got1 := p.snippetStart(tt.args.line)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("snippetStart() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("snippetStart() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_snippet_String(t *testing.T) {
	type fields struct {
		meta    metadata
		content []string
	}
	type args struct {
		targetIndent string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			fields: fields{
				content: []string{"\tfoo", "\t\tbar"},
				meta: metadata{
					sourceIndent: "\t",
					tag:          "",
				},
			},
			args: args{targetIndent: " "},
			want: " foo\n \tbar",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := snippet{
				meta:    tt.fields.meta,
				content: tt.fields.content,
			}
			if got := s.String(tt.args.targetIndent); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getSnippetsFromContent(t *testing.T) {
	goParser, _ := newParser(".go")

	type args struct {
		p       parser
		content string
	}
	tests := []struct {
		name    string
		args    args
		want    snippets
		wantErr bool
	}{
		{
			name: "no indentation",
			args: args{
				p: goParser,
				content: `
// START SNIPPET snippet1
func main() {
	continue
}
// END SNIPPET
`,
			},
			want: snippets{
				"snippet1": snippet{
					meta: metadata{
						sourceIndent: "",
						tag:          "snippet1",
					},
					content: []string{"func main() {", "\tcontinue", "}"},
				},
			},
			wantErr: false,
		}, {
			name: "indentation",
			args: args{
				p: goParser,
				content: `
	// START SNIPPET snippet1
	func main() {
		continue
	}
	// END SNIPPET
	// START SNIPPET snippet2
	some
		other
	content
	// END SNIPPET
`,
			},
			want: snippets{
				"snippet1": snippet{
					meta: metadata{
						sourceIndent: "\t",
						tag:          "snippet1",
					},
					content: []string{"\tfunc main() {", "\t\tcontinue", "\t}"},
				},
				"snippet2": snippet{
					meta: metadata{
						sourceIndent: "\t",
						tag:          "snippet2",
					},
					content: []string{"\tsome", "\t\tother", "\tcontent"},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := snippetsFromContent(tt.args.p, tt.args.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("snippetsFromContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("snippetsFromContent() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_snippets_addSnippets(t *testing.T) {
	s := snippets{"foo": snippet{content: []string{"bar"}}}
	ext := snippets{"bar": snippet{content: []string{"baz"}}}

	want := snippets{
		"foo": snippet{content: []string{"bar"}},
		"bar": snippet{content: []string{"baz"}},
	}

	s.addSnippets(ext)

	if !reflect.DeepEqual(s, want) {
		t.Errorf("got: %v, want: %v", s, want)
	}
}

func Test_parser_snippetLocation1(t *testing.T) {
	goParser, _ := newParser(".go")

	type args struct {
		line string
	}
	tests := []struct {
		name  string
		p     parser
		args  args
		want  string
		want1 string
		want2 bool
	}{
		{
			p:     goParser,
			args:  args{line: "  // PUT SNIPPET foobar"},
			want:  "  ",
			want1: "foobar",
			want2: true,
		}, {
			name:  "white-space at the end",
			p:     goParser,
			args:  args{line: "  // PUT SNIPPET foobar\t "},
			want:  "  ",
			want1: "foobar",
			want2: true,
		}, {
			name:  "unicode-space at the end",
			p:     goParser,
			args:  args{line: "  // PUT SNIPPET foobar　"},
			want:  "",
			want1: "",
			want2: false,
		}, {
			p:     goParser,
			args:  args{line: "  // END SNIPPET"},
			want:  "",
			want1: "",
			want2: false,
		}, {
			p:     goParser,
			args:  args{line: "some content"},
			want:  "",
			want1: "",
			want2: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.p
			got, got1, got2 := p.snippetLocation(tt.args.line)
			if got != tt.want {
				t.Errorf("snippetLocation() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("snippetLocation() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("snippetLocation() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}
