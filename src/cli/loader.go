package cli

import "github.com/maasasia/donggu/exporter"

var fileExporters = map[string]exporter.DictionaryFileExporter{
	"json": exporter.JsonDictionaryExporter{},
}

var projectExporters = map[string]exporter.DictionaryProjectExporter{
	"typescript": exporter.TypescriptDictionaryExporter{},
}

func loadProjectExporter(name string) exporter.DictionaryProjectExporter {
	projectExporter, isProjectExporter := projectExporters[name]

	if isProjectExporter {
		return projectExporter
	}
	return nil
}

func loadFileExporter(name string) exporter.DictionaryFileExporter {
	fileExporter, isFileExporter := fileExporters[name]

	if isFileExporter {
		return fileExporter
	}
	return nil
}
