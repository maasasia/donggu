package util

import (
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

type ReplaceSet struct {
	From *regexp.Regexp
	To   string
}

// Warning: This loads the full content of each file into memory.
func BatchReplaceFiles(files []string, from, to string) error {
	for _, file := range files {
		fileContent, err := ioutil.ReadFile(file)
		if err != nil {
			return errors.Wrapf(err, "failed to open '%s'", file)
		}
		fileContent = []byte(strings.ReplaceAll(string(fileContent), from, to))
		if err := ioutil.WriteFile(file, fileContent, os.ModePerm); err != nil {
			return errors.Wrapf(err, "failed to write '%s'", file)
		}
	}
	return nil
}

func MultiReplaceFile(file string, replaces []ReplaceSet) error {
	fileBinary, err := ioutil.ReadFile(file)
	if err != nil {
		return errors.Wrapf(err, "failed to open '%s'", file)
	}
	fileContent := string(fileBinary)
	for _, replace := range replaces {
		fileContent = replace.From.ReplaceAllString(fileContent, replace.To)
	}
	if err := ioutil.WriteFile(file, []byte(fileContent), os.ModePerm); err != nil {
		return errors.Wrapf(err, "failed to write '%s'", file)
	}
	return nil
}
