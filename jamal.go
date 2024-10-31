package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"strings"

	"github.com/docopt/docopt-go"
	"github.com/ghodss/yaml"
)

var (
	err         error
	data        []byte
	decodedData []byte
	logger      = log.New(os.Stderr, "", 0)
	y2jOptions  = []string{"y2j", "yaml2json", "yamltojson"}
	j2yOptions  = []string{"j2y", "json2yaml", "jsontoyaml"}
)

func isFound(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func yamlToJson(raw []byte) ([]byte, error) {
	var (
		data   interface{}
		output []byte
	)

	err := yaml.Unmarshal(raw, &data)

	if err != nil {
		return nil, err
	}

	output, err = json.MarshalIndent(data, "", "  ")
	output = append(output, "\n"...)

	return output, err
}

// Converts JSON to YAML.
func jsonToYaml(raw []byte) ([]byte, error) {
	var data interface{}

	err := json.Unmarshal(raw, &data)

	if err != nil {
		return nil, err
	}

	output, err := yaml.Marshal(data)

	return output, err
}

func readStdin() ([]byte, error) {
	fi, err := os.Stdin.Stat()

	if err != nil {
		return nil, err
	}

	if fi.Mode()&os.ModeNamedPipe == 0 {
		err = errors.New("nothing piped into stdin")

		return nil, err
	}

	return io.ReadAll(os.Stdin)
}

func logError(msg string, err error) {
	logger.Println(msg)
	logger.Println(err.Error())
}

func main() {
	doc := `Jamal

        Command line interface for converting JSON to YAML and YAML to JSON.
        Expects either an input file or data from stdin.

        Usage:
            jamal <action> [<input-file>]

        Options:
            -h --help       Show this message.

        Arguments:
            <action>        Conversion action.
                            [yamltojson, y2j, yaml2json | jsontoyaml, j2y, json2yaml]

            <input-file>    Path to data file.
    `
	arguments, _ := docopt.Parse(doc, nil, true, "Jamal 1.0.0", false)
	dataPath := arguments["<input-file>"]
	action := arguments["<action>"].(string)
	action = strings.ToLower(action)

	if dataPath == nil {
		data, err = readStdin()
	} else {
		path := dataPath.(string)
		data, err = os.ReadFile(path)
	}

	if err != nil {
		logError("Error occurred loading data.", err)

		os.Exit(1)
	}

	switch {
	case isFound(y2jOptions, action):
		decodedData, err = yamlToJson(data)
	case isFound(j2yOptions, action):
		decodedData, err = jsonToYaml(data)
	default:
		logger.Println("Invalid action.")
		logger.Println(doc)
		os.Exit(1)
	}

	if err != nil {
		logError("Error occurred converting data.", err)

		os.Exit(1)
	}

	os.Stdout.Write(decodedData)
}
