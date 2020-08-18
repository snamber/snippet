package snippet

import "testing"

func Test_replace(t *testing.T) {
	goParser, _ := newParser(".go")
	type args struct {
		p       parser
		content string
		s       snippets
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{
				p: goParser,
				content: `first line
 // PUT SNIPPET foobar
last line
`, // mind the space in front of '// SNIPPET'
				s: snippets{"foobar": snippet{
					meta: metadata{
						sourceIndent: "",
						tag:          "foobar",
					},
					content: []string{"\tindented", "\tmultiline", "\tcontent"},
				}},
			},
			want: `first line
 	indented
 	multiline
 	content
last line
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := replace(tt.args.p, tt.args.content, tt.args.s); got != tt.want {
				t.Errorf("replace() = %v, want %v", got, tt.want)
			}
		})
	}
}
