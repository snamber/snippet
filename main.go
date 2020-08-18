package snippet

import "github.com/pkg/errors"

func ExtractAndExpandSnippets(snippetsPath, templatePath, outputPath string) error {
	s, err := extractSnippetsFromDir(snippetsPath)
	if err != nil {
		return errors.Wrap(err, "could not extract snippets")
	}

	err = expandTemplates(templatePath, outputPath, s)
	if err != nil {
		return errors.Wrap(err, "could not expand templates")
	}

	return nil
}
