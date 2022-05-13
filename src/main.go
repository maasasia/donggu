package main

import "github.com/maasasia/donggu/cli"

func main() {
	cli.Execute()

	// fmt.Println(code.CopyTemplateTo("typescript", "/Users/mac/Desktop/template_test", code.CopyTemplateOptions{}))
	// im := importer.JsonDictionaryImporter{}
	// var cim importer.DictionaryFileImporter = importer.JsonDictionaryImporter{}
	// pj, err := importer.JsonDictionaryImporter{}.OpenMetadataFile("/Users/mac/Desktop/dgtest")
	// if err != nil {
	// 	panic(err)
	// }
	// mt, err := im.ImportMetadata(pj)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(mt)
	// file, err := importer.JsonDictionaryImporter{}.OpenContentFile("/Users/mac/Desktop/dgtest")
	// if err != nil {
	// 	panic(err)
	// }
	// defer file.Close()

	// ct, err := cim.ImportContent(file, mt)
	// if err != nil {
	// 	panic(err)
	// }
	// ct.ToTree().Print()

	// var encoder exporter.DictionaryProjectExporter = exporter.TypescriptDictionaryExporter{}
	// err = encoder.Export("generated", ct, mt)
	// fmt.Println(err)
}
