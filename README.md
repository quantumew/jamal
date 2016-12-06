Jamal
=====

YAML to JSON and JSON to YAML conversion tool. Written in Golang so it can be statically compiled and dropped anywhere you need it.

[Build of latest release.](https://github.com/quantumew/jamal/releases)

Build
-----

    go get github.com/quantumew/jamal
    cd "$GOPATH/src/github.com/quantumew/jamal"
    go build
    mv jamal <in you path somewhere>

Usage
-----

    ./jamal <action> [<input-file>]

    Options:
        -h --help       Show this message.

    Arguments:
        <action>        Conversion action. [yamltojson, y2j, yaml2json | jsontoyaml, j2y, json2yaml]

        <input-file>    Path to data file.


Examples
--------

    # Output YAML from a JSON file.
    jamal json2yaml some-file.json

    # Output JSON from a YAML file.
    jamal yaml2json some-file.yaml

    # Output YAML from JSON from STDIN.
    echo '{"property": 5, "otherThing": true}' | jamal j2y

