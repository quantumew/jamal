//usr/bin/env go run $0 $@; exit;
package main

import (
	"encoding/json"
	"errors"
	"github.com/docopt/docopt-go"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var logger = log.New(os.Stderr, "", 0)

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

	var (
		err         error
		data        []byte
		decodedData []byte
	)

	// Sort of ugly but this version of docopt does not support this
	// type of validation.
	if action != "yaml2json" && action != "yamltojson" && action != "y2j" {
		if action != "json2yaml" && action != "jsontoyaml" && action != "j2y" {
			logger.Println("Invalid action.")
			logger.Println(doc)

			os.Exit(1)
		}
	}

	if dataPath == nil {
		data, err = readStdin()
	} else {
		path := dataPath.(string)
		data, err = ioutil.ReadFile(path)
	}

	if err != nil {
		logError("Error occurred loading data.", err)

		os.Exit(1)
	}

	if action == "yaml2json" || action == "yamltojson" || action == "y2j" {
		decodedData, err = yamlToJson(data)
	} else {
		decodedData, err = jsonToYaml(data)
	}

	if err != nil {
		logError("Error occurred converting data.", err)

		os.Exit(1)
	}

	os.Stdout.Write(decodedData)
}

// Converts YAML to JSON.
func yamlToJson(raw []byte) ([]byte, error) {
	var data interface{}

	err := yaml.Unmarshal(raw, &data)

	if err != nil {
		return nil, err
	}

	output, err := json.MarshalIndent(data, "", "  ")

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
		err = errors.New("Nothing piped into stdin.")

		return nil, err
	}

	return ioutil.ReadAll(os.Stdin)
}

func logError(msg string, err error) {
	logger.Println(msg)
	logger.Println(err.Error())
}
