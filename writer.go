package snippet

import (
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// expandTemplates expands the templates in src with snippets s
func expandTemplates(sourceDir, destDir string, s snippets) error {

	f := func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		// pass directories
		if info.IsDir() {
			return nil
		}

		// read file
		bt, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		// expand snippets
		parser, err := newParser(filepath.Ext(path))
		if err != nil {
			log.Println(errors.Wrapf(err, "warning: could not expand snippets in file %s", path))
			return nil
		}
		out := replace(parser, string(bt), s)

		// write content
		rel, err := filepath.Rel(filepath.Dir(sourceDir), filepath.Dir(path))
		if err != nil {
			return err
		}
		newPath := filepath.Clean(filepath.Join(destDir, rel, info.Name()))
		err = ioutil.WriteFile(newPath, []byte(out), info.Mode())
		if err != nil {
			return err
		}
		return nil
	}

	err := filepath.Walk(sourceDir, f)
	if err != nil {
		return errors.Wrapf(err, "could not successfully expand templates for dir %s", sourceDir)
	}

	return nil
}
