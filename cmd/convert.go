package cmd

import (
	"fmt"
	"os"

	"github.com/mheers/raygun2x/openapi"
	"github.com/mheers/raygun2x/postman"
	"github.com/mheers/raygun2x/raygun"
	"github.com/spf13/cobra"
)

var (
	inputFolder string
	outputFile  string
	mode        string
	stdout      bool

	convertCmd = &cobra.Command{
		Use:   "convert",
		Short: "Convert test suites to OpenAPI or Postman Collection",
		Run: func(cmd *cobra.Command, args []string) {
			convert()
		},
	}
)

func init() {
	convertCmd.Flags().StringVarP(&inputFolder, "input", "i", "presets/", "Input folder containing test suites")
	convertCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file (if empty, prints to stdout)")
	convertCmd.Flags().StringVarP(&mode, "mode", "m", "openapi", "Conversion mode: 'openapi' or 'postman'")
	convertCmd.Flags().BoolVarP(&stdout, "stdout", "s", false, "Print output to stdout instead of file")
}

func convert() {
	suites, err := raygun.GetTestSuites(inputFolder)
	if err != nil {
		panic(err)
	}

	var result []byte

	switch mode {
	case "postman":
		collection, err := postman.CreateCollection(suites)
		if err != nil {
			panic(err)
		}
		result, err = postman.MarshalCollection(collection)
		if err != nil {
			panic(err)
		}
	case "openapi":
		spec, err := openapi.CreateOpenAPISpec(suites)
		if err != nil {
			panic(err)
		}
		result, err = openapi.MarshalSpec(spec)
		if err != nil {
			panic(err)
		}
	default:
		panic("Invalid mode. Use 'openapi' or 'postman'")
	}

	if stdout || outputFile == "" {
		fmt.Println(string(result))
	} else {
		if err := writeToFile(outputFile, result); err != nil {
			panic(err)
		}
	}
}

func writeToFile(filename string, data []byte) error {
	return os.WriteFile(filename, data, 0644)
}
