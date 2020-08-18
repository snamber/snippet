package snippet

import (
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"testing"
)

func Test_expandTemplates(t *testing.T) {
	s, err := extractSnippetsFromDir(path.Join("snippet", "test"))
	if err != nil {
		t.Error(errors.Wrap(err, "can't extract example snippets for test"))
	}
	tmpDir, err := ioutil.TempDir("", "snippet-master")
	defer os.RemoveAll(tmpDir)

	type args struct {
		sourceDir string
		destDir   string
		s         snippets
	}
	_ = s
	test := struct {
		name    string
		args    args
		wantErr bool
	}{
		args: args{
			sourceDir: path.Join("test", "input"),
			destDir:   tmpDir,
			s:         s,
		},
		wantErr: false,
	}

	t.Run(test.name, func(t *testing.T) {
		if err := expandTemplates(test.args.sourceDir, test.args.destDir, test.args.s); (err != nil) != test.wantErr {
			t.Errorf("expandTemplates() error = %v, wantErr %v", err, test.wantErr)
		} else {
			// check if template file is the same
			gotFile := path.Join(tmpDir, "readme.md")
			_, err := os.Stat(gotFile)
			if err != nil {
				t.Error(errors.Wrap(err, "expected output file does not exist"))
			}

			gotBt, err := ioutil.ReadFile(gotFile)
			if err != nil {
				t.Error(errors.Wrap(err, "could not read output file"))
			}
			wantBt, err := ioutil.ReadFile(path.Join("test", "output", "readme.md"))
			if err != nil {
				t.Error(errors.Wrap(err, "could not read target file"))
			}

			if reflect.DeepEqual(gotBt, wantBt) {
				t.Errorf("got: \n%s\n wanted: \n%s\n", string(gotBt), string(wantBt))
			}
		}
	})
}
