package cli

import (
	"github.com/maasasia/donggu/exporter"
	"github.com/maasasia/donggu/importer"
)

var fileImporters = map[string]importer.DictionaryFileImporter{
	"json": importer.JsonDictionaryImporter{},
	"csv":  importer.CsvDictionaryImporter{},
}

var importers = map[string]importer.DictionaryImporter{
	"json": importer.JsonDictionaryImporter{},
	"csv":  importer.CsvDictionaryImporter{},
}

var fileExporters = map[string]exporter.DictionaryFileExporter{
	"json": exporter.JsonDictionaryExporter{},
	"csv":  exporter.CsvDictionaryExporter{},
}

var projectExporters = map[string]exporter.DictionaryProjectExporter{
	"typescript": exporter.TypescriptDictionaryExporter{},
	"golang":     exporter.GolangDictionaryExporter{},
}

func loadImporter(name string) importer.DictionaryImporter {
	importer, isImporter := importers[name]

	if isImporter {
		return importer
	}
	return nil
}

func loadFileImporter(name string) importer.DictionaryFileImporter {
	importer, isImporter := fileImporters[name]
	if isImporter {
		return importer
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

func loadProjectExporter(name string) exporter.DictionaryProjectExporter {
	projectExporter, isProjectExporter := projectExporters[name]

	if isProjectExporter {
		return projectExporter
	}
	return nil
}
