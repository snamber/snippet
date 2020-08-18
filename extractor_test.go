package snippet

import (
	"path"
	"reflect"
	"testing"
)

func Test_extractSnippetsFromDir(t *testing.T) {
	type args struct {
		root string
	}
	tests := []struct {
		name string
		args args
		want snippets
	}{
		{
			args: args{root: path.Join("test")},
			want: snippets{
				"someFunc": snippet{
					meta: metadata{
						sourceIndent: "\t",
						tag:          "someFunc",
					},
					content: []string{"\tl := someFunc(\"some input string\")", "\t// l == 17"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := extractSnippetsFromDir(tt.args.root); !reflect.DeepEqual(got, tt.want) {
				if err != nil {
					t.Error(err)
				}
				t.Errorf("extractSnippetsFromDir() = %v, want %v", got, tt.want)
			}
		})
	}
}
