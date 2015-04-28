// Copyright (c) 2013-2015 by The SeedStack authors. All rights reserved.

// This file is part of SeedStack, An enterprise-oriented full development stack.

// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"
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
var dirPath = "./"

func init() {
	flag.StringVar(&transPath, "t", "./tdf.yml", "Specify the path to the transformation description file")
	flag.BoolVar(&verbose, "v", false, "Enable verbose mode.")
	flag.Parse()
}

func main() {
	if flag.Arg(0) == "fix" {
		fix()
	} else if flag.Arg(0) == "help" {
		if flag.Arg(1) == "fix" {
			fmt.Println(`Fix the files in a given directory based on a YAML transformation description file. 
If no directory is passed as argument, the transformations will be applied on current directory.

Usage: 
  seed [flags] fix [directory/to/transform]

Available flags:
 -t file/path.yml: the YAML transformation description file

YAML transformation description file format:

The description file accepts a list of transformation. Each transformation can have include files or exclude directories. 
It can also use higher level preconditions with "pre" which uses the file content. Finally, it takes a list of procedure to apply the file. 
Procedures are described with their name and the arguments to pass. See the following 'tdf.yaml' file as example. 

tdf.yml
----------------
- 
 Include: "*.go|*.yml"
 Exclude: "*.out"
 pre: 
  - AlwaysTrue
  - ...
 proc:
  - Replace
   Name: Replace
   Params:
    - "old"
    - "new"
  - ...
-
 ...
----------------
`)
		}
	} else {
		fmt.Println(`Usage: seed <command> <args>

Commands:
    fix    Apply source transformation on a directory, based on a YAML transformation file
    help   Provide help for seed commands 

See 'seed help <command>' to read about a specific subcommand.
`)
	}
}

func fix() {
	start := time.Now()

	var dat []byte
	var tdfPath string

	if transPath != strings.TrimPrefix(transPath, "http://") ||
		transPath != strings.TrimPrefix(transPath, "https://") {
		// get the transformation description file from internet

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

		dat = body
	} else {
		// get the tdf from the file system
		absPath, errFilePath := filepath.Abs(transPath)
		tdfPath = absPath

		if errFilePath != nil {
			log.Fatal("Error constructing the file path.\n", errFilePath)
		}

		bytes, err := ioutil.ReadFile(tdfPath)
		if err != nil {
			log.Fatal("Unable to read the transformation description file.\n", err)
		}
		dat = bytes
	}

	if verbose {
		fmt.Printf("Parse the transformation description file: %s.\n", transPath)
	}
	transf := parseTdf(dat)

	// set the directory to parse if specified
	if flag.Arg(1) != "" {
		absPath, errFilePath := filepath.Abs(flag.Arg(1))
		if errFilePath != nil {
			log.Fatal("Error constructing the file path.\n", errFilePath)
		}
		dirPath = absPath
	}

	count := processFiles(walkDir(dirPath, transf.Exclude, tdfPath), transf)

	elapsed := time.Since(start)
	fmt.Printf("%s %s fixed: %v files\n", filepath.Base(dirPath), elapsed, count)
}

func parseTdf(dat []byte) T {
	var t T
	err := yaml.Unmarshal(dat, &t)
	if err != nil {
		log.Fatal("Failed to parse the transformation description file.\n", err)
	}
	return t
}
