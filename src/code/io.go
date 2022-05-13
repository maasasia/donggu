package code

import (
	"os"
	"path"

	cp "github.com/otiai10/copy"
	"github.com/pkg/errors"
)

type CopyTemplateOptions struct {
	Skip func(src string) (bool, error)
}

func CopyTemplateTo(templateName, destination string, options CopyTemplateOptions) error {
	if err := os.MkdirAll(destination, os.ModePerm); err != nil {
		return errors.Wrapf(err, "failed to copy template '%s'", templateName)
	}
	source, err := resolveTemplateFolder(templateName)
	if err != nil {
		return errors.Wrapf(err, "failed to copy template '%s'", templateName)
	}
	err = cp.Copy(source, destination, cp.Options{Skip: options.Skip})
	return errors.Wrapf(err, "failed to copy template '%s'", templateName)
}

func resolveTemplateFolder(templateName string) (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", errors.Wrap(err, "failed to resolve template folder")
	}
	return path.Join(execPath, "..", "templates", templateName), nil
}
