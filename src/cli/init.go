package cli

import (
	"os"
	"sort"
	"strings"

	"github.com/maasasia/donggu/dictionary"
	"github.com/maasasia/donggu/exporter"
	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func runInitCommand(cmd *cobra.Command, _ []string) error {
	projectRoot, _ := cmd.Flags().GetString("project")
	if projectRoot == "" {
		prompt := promptui.Prompt{
			Label: "Directory to initialize? (leave empty for current directory)",
		}
		result, err := prompt.Run()
		if err != nil {
			return errors.Wrap(err, "prompt failed")
		}
		if result == "" {
			workingDir, err := os.Getwd()
			if err != nil {
				return errors.Wrap(err, "failed to resolve cwd")
			}
			projectRoot = workingDir
		} else {
			projectRoot = result
		}
	}

	versionPrompt := promptui.Prompt{Label: "Input initial version (leave empty for 0.0.1)"}
	version, err := versionPrompt.Run()
	if err != nil {
		return errors.Wrap(err, "prompt failed")
	}
	version = strings.TrimSpace(version)
	if version == "" {
		version = "0.0.1"
	}

	allLangPrompt := promptui.Prompt{Label: "Input all supported languages, seperated by a comma. (ex: ko, en, ja)"}
	allLang, err := allLangPrompt.Run()
	if err != nil {
		return errors.Wrap(err, "prompt failed")
	}
	supportedLanguageList := strings.Split(allLang, ",")
	for index := range supportedLanguageList {
		supportedLanguageList[index] = strings.TrimSpace(supportedLanguageList[index])
		if supportedLanguageList[index] == "" {
			return errors.Wrap(err, "invalid input: found empty language")
		}
	}

	requiredLangPrompt := promptui.Prompt{Label: "Input all languages that are required (ie. languages that all entries should support)"}
	requiredLang, err := requiredLangPrompt.Run()
	if err != nil {
		return errors.Wrap(err, "prompt failed")
	}
	requiredLangList := strings.Split(requiredLang, ",")
	for index := range requiredLangList {
		requiredLangList[index] = strings.TrimSpace(requiredLangList[index])
		if requiredLangList[index] == "" {
			return errors.Wrap(err, "invalid input: found empty language")
		}
	}

	sort.Strings(requiredLangList)
	sort.Strings(supportedLanguageList)

	metadata := dictionary.Metadata{
		Version:            version,
		RequiredLanguages:  requiredLangList,
		SupportedLanguages: supportedLanguageList,
		ExporterOptions:    map[string]map[string]interface{}{},
	}

	if err := metadata.Validate(); err != nil {
		return errors.Wrap(err, "invalid input")
	}

	content := dictionary.FlattenedContent{}
	content["example"] = dictionary.Entry{}
	requiredLangSet := metadata.RequiredLanguageSet()
	for lang := range metadata.SupportedLanguageSet() {
		if _, ok := requiredLangSet[lang]; ok {
			content["example"][lang] = "This language is required."
		} else {
			content["example"][lang] = "This language is optional, so this key may be deleted."
		}
	}

	err = exporter.JsonDictionaryExporter{}.Export(
		projectRoot,
		&content,
		metadata,
		exporter.OptionMap{},
	)
	if err != nil {
		return errors.Wrap(err, "failed to write project")
	}
	return nil
}

func initInitCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize new project",
		Run:   wrapExecCommand(runInitCommand),
	}
	return cmd
}
