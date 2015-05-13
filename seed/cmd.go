// Copyright (c) 2013-2015 by The SeedStack authors. All rights reserved.

// This file is part of SeedStack, An enterprise-oriented full development stack.

// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"
	"os"
	"bufio"
)

const (
	fixHelp = `Usage: seed [flags] fix [directory/to/transform]

Fix the files in a given directory, based on a transformation file. If no directory 
is passed as argument, the transformations will be applied in current directory.

Available flags:
 -t file/path.yml  the YAML transformation file
 -v                verbose mode
 -vv               very verbose mode

Transformation file:

The description file can be in YAML or TOML. It accepts a list of transformations.
Transformation contains filter based on the file name, but they can also use higher
order precondition based on the file content. Finally, it takes a list of procedures
to apply on the file. Procedures are described with a name and a list of arguments.
This files also accepts global exclusions based on directory names. The directories 
to exclude are separated by "|".

Sample of YAML transformation file.

----------------
exclude: "target|.git"
transformations:
 -
  filter: "pom.xml"
  pre: 
    - AlwaysTrue
    - ...
  proc:
    -
      name: Replace
      params:
        - "old"
        - "new"
    - 
      ...
 -
  ...
----------------

A convert method exists to convert yaml into toml see "seed convert [file] [format]".
`
	seedHelp = `Usage: seed <command> <args>

Commands:
    fix      Apply source transformations on a directory
    convert  Convert a yaml transformation file into toml
    help     Provide help for seed commands 
    version  Show the seed tool version

See 'seed help <command>' to read about a specific subcommand.
`
)

// T correspond to the content of a transformation file.
// It contains exclude directories and an array of transformations.
type T struct {
	Exclude         string
	Transformations []Transformation
}

// Transformation is a strutucture representating a set
// of procedure to apply on a source code directory
type Transformation struct {
	Filter string
	Pre    []string
	Proc   []Procedure
}

// Procedure is a function call with a method name and
// its parameters
type Procedure struct {
	Name   string
	Params []string
}

var transPath string
var verbose bool
var vverbose bool
var dirPath = "./"

func init() {
	flag.StringVar(&transPath, "t", "./tdf.yml", "Specify the path to the transformation description file")
	flag.BoolVar(&verbose, "v", false, "Enable verbose mode.")
	flag.BoolVar(&vverbose, "vv", false, "Enable very verbose mode.")
	flag.Parse()

	if vverbose {
		verbose = true
	}
}

func main() {
	switch flag.Arg(0) {
	case "fix":
		fix()
	case "convert":
		convertTdf(flag.Arg(1), flag.Arg(2))
	case "help":
		if flag.Arg(1) == "fix" {
			fmt.Println(fixHelp)
		}
	case "version":
		fmt.Println("Seed Tool v0.1")
	default:
		fmt.Println(seedHelp)
	}
}

func fix() {
	start := time.Now()

	var dat []byte
	var tdfPath string

	if strings.HasPrefix(transPath, "http://") || strings.HasPrefix(transPath, "https://") {
		dat = fetchURL(transPath)
	} else {
		dat = readFile(transPath)
	}

	if verbose {
		fmt.Printf("Apply transformations from: %s.\n\n---\n", transPath)
	}

	format, err := getFormat(transPath)
	if err != nil {
		log.Fatalf("Unsupported format for %s", transPath)
	}
	transf := parseTdf(dat, format)

	// set the directory to parse if specified
	if flag.Arg(1) != "" {
		absPath, errFilePath := filepath.Abs(flag.Arg(1))
		if errFilePath != nil {
			log.Fatal("Error constructing the file path.\n", errFilePath)
		}
		dirPath = absPath
	}

	files := walkDir(dirPath, transf.Exclude, tdfPath)
	count := processFiles(files, transf)

	elapsed := time.Since(start)
	var shortDirPath = filepath.Base(dirPath)
	if shortDirPath == "." {
		wd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get current dir: %s", err)
		}
		shortDirPath = filepath.Base(wd)
	}
	fmt.Printf("\n%s fixed %v/%v files in %s\n", shortDirPath, count, len(files), elapsed)
}

func getFormat(name string) (string, error) {
	index := strings.LastIndex(name, ".") + 1
	extension := strings.ToLower(name[index:])

	var ext string
	var err error

	switch extension {
	case "yml", "yaml":
		ext = "yml"
	case "toml":
		ext = "toml"
	default:
		err = fmt.Errorf("%s format unsupported", extension)
	}

	return ext, err
}

func fetchURL(url string) []byte {
	resp, err := http.Get(transPath)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode > 299 {
		log.Fatalf("Error %v when fetching %s\n", resp.StatusCode, transPath)
	}

	body, err2 := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatal("Error reading http reponse.\n", err2)
	}

	return body
}

func readFile(path string) []byte {
	absPath, errFilePath := filepath.Abs(path)
	tdfPath := absPath

	if errFilePath != nil {
		log.Fatal("Error constructing the file path.\n", errFilePath)
	}

	bytes, err := ioutil.ReadFile(tdfPath)
	if err != nil {
		log.Fatal("Unable to read the transformation description file.\n", err)
	}
	return bytes
}

func parseTdf(dat []byte, format string) T {
	var t T

	switch format {
	case "yml":
		err := yaml.Unmarshal(dat, &t)
		if err != nil {
			log.Fatalf("Failed to parse the yaml file: %s", err)
		}
	case "toml":
		err := toml.Unmarshal(dat, &t)
		if err != nil {
			log.Fatalf("Failed to parse the toml file: %s", err)
		}
	}
	return t
}

func convertTdf(path, newFormat string) {
	index := strings.LastIndex(path, ".") + 1
	format, err := getFormat(path[index:])
	if err != nil {
		log.Fatalf("Unsupported format for %s", path)
	}
	
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Unable to read the transformation description file.\n", err)
	}
	var t T
	var res []byte

	f, err := os.Create(path[:index] + "toml")
	if err != nil {
		log.Fatal(err)
	}
	
	writer := bufio.NewWriter(f)
	
	if newFormat == "toml" && (format == "yaml" || format == "yml") {
		err := yaml.Unmarshal(bytes, &t)
		if err != nil {
			log.Fatalf("Failed to parse the yaml file: %s", err)
		}

		if err = toml.NewEncoder(writer).Encode(t); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatalf("%s format is not supported for convertion", format)
	}
	
	ioutil.WriteFile(path[:index], res, 0666)
}
