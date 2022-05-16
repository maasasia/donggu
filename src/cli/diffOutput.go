package cli

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/maasasia/donggu/dictionary"
	"github.com/pkg/errors"
	"github.com/rodaine/table"
)

func outputDiffToConsole(diff dictionary.ContentDifference, reverse bool) {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	createdColorFmt := color.New(color.FgBlue).SprintFunc()
	deletedColorFmt := color.New(color.FgRed).SprintfFunc()
	changedColorFmt := color.New(color.FgYellow).SprintFunc()

	tbl := table.New("Key", "Difference").WithHeaderFormatter(headerFmt)
	diffMap := map[string]string{}

	if reverse {
		fmt.Println("Changes from given file to project file (content.json)")
	} else {
		fmt.Println("Changes from project file (content.json) to given file")
	}

	fmt.Printf("- Number of keys in source: %d\n", diff.SourceKeyCount)
	fmt.Printf("- Number of keys in destination: %d\n", diff.DestKeyCount)
	fmt.Printf("- Number of unchanged keys: %d\n", diff.UnchangedKeyCount)

	for _, key := range diff.CreatedKeys {
		diffMap[string(key)] = createdColorFmt("Created")
	}
	for _, key := range diff.DeletedKeys {
		diffMap[string(key)] = deletedColorFmt("Deleted")
	}

	for key, value := range diff.Changes {
		changes := []string{}

		for lang, change := range value {
			changes = append(changes, fmt.Sprintf("%s: %s", lang, change))
		}
		sort.Strings(changes)
		diffMap[string(key)] = fmt.Sprintf(
			"%s(%s)",
			changedColorFmt("Changed"),
			strings.Join(changes, ", "),
		)
	}

	changeKeys := make([]string, 0, len(diffMap))
	for key := range diffMap {
		changeKeys = append(changeKeys, key)
	}
	sort.Strings(changeKeys)

	for _, key := range changeKeys {
		tbl.AddRow(key, diffMap[key])
	}
	tbl.Print()
}

func outputDiffToCsv(diff dictionary.ContentDifference) error {
	diffMap := map[string][]string{}

	for _, key := range diff.CreatedKeys {
		diffMap[string(key)] = []string{"Created", ""}
	}
	for _, key := range diff.DeletedKeys {
		diffMap[string(key)] = []string{"Deleted", ""}
	}
	for key, value := range diff.Changes {
		changes := []string{}

		for lang, change := range value {
			changes = append(changes, fmt.Sprintf("%s: %s", lang, change))
		}
		sort.Strings(changes)
		diffMap[string(key)] = []string{
			"Changed",
			strings.Join(changes, ", "),
		}
	}

	changeKeys := make([]string, 0, len(diffMap))
	for key := range diffMap {
		changeKeys = append(changeKeys, key)
	}
	sort.Strings(changeKeys)

	file, err := os.OpenFile("diff.csv", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "failed to open csv file")
	}
	defer file.Close()
	csvWriter := csv.NewWriter(file)
	defer csvWriter.Flush()
	csvWriter.Write([]string{"Key", "Difference", "Details"})

	for _, key := range changeKeys {
		csvWriter.Write([]string{key, diffMap[key][0], diffMap[key][1]})
	}

	return nil
}
