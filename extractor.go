package snippet

import (
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func extractSnippetsFromDir(root string) (snippets, error) {
	ss := snippets{}

	walkfunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}

		bt, err := ioutil.ReadFile(path)
		if err != nil {
			log.Printf("skipping %s: %s", path, err)
			return nil
		}

		ext := filepath.Ext(info.Name())
		p, err := newParser(ext)
		if err != nil {
			log.Printf("skipping %s: %s", path, err)
			return nil
		}

		s, err := snippetsFromContent(p, string(bt))
		if err != nil {
			log.Printf("skipping %s: %s", path, err)
			return nil
		}

		ss.addSnippets(s)
		return nil
	}

	err := filepath.Walk(root, walkfunc)
	if err != nil {
		return snippets{}, errors.Wrapf(err, "failed extracting snippets from %v", root)
	}
	return ss, nil
}
