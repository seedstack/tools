// Copyright (c) 2013-2015 by The SeedStack authors. All rights reserved.

// This file is part of SeedStack, An enterprise-oriented full development stack.

// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"path/filepath"
	"flag"
	"io/ioutil"
	"strings"
	"net/http"
	"log"
)

// Transformation is a strutucture representating a set 
// of procedure to apply on a source code directory
type Transformation struct {
	Pre []string
	Proc []Procedure
}

// Procedure is a function call with a method name and 
// its parameters
type Procedure struct {
	Name string
	Params []string
}

var transPath string
var dirPath = "./"

func init() {
	flag.StringVar(&transPath, "t", "./tdf.yml", "Specify the path to the transformation description file")
	
	flag.Parse()
}

func main() {
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
			log.Fatalf("Error %v when fetching %s", resp.StatusCode, transPath)
		}

		body, err2 := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Fatal(err2)
		}

		dat = body
	} else {
		// get the tdf from the file system

		tdfPath, errFilePath := filepath.Abs(transPath)
		if errFilePath != nil {
			log.Fatal("Error constructing the file path.", errFilePath)
		}
		
		bytes, err := ioutil.ReadFile(tdfPath)
		if err != nil {
			log.Fatal(err)
		}
		dat = bytes
	}

	transf := parseTdf(dat)
	fmt.Printf("Parsed the transformation description file from %s.\n", transPath)

	// set the directory to parse if specified
    if flag.Arg(0) != "" {
		absPath, errFilePath := filepath.Abs(flag.Arg(0))
		if errFilePath != nil {
			log.Fatal("Error constructing the file path.", errFilePath)
		}
		dirPath = absPath
	}

	processFiles(walkthroughDir(tdfPath, dirPath), transf)
}

func parseTdf(dat []byte) []Transformation {
	var transf []Transformation
	err := yaml.Unmarshal(dat, &transf)
	if err != nil {
		log.Fatal("Failed to parse the transformation description file.\n", err)
	}

	return transf
}
