// Copyright (c) 2013-2015 by The SeedStack authors. All rights reserved.

// This file is part of SeedStack, An enterprise-oriented full development stack.

// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func walkDir(root string, excludes string, tdfPath string) []string {
	var files []string
	if vverbose {
		fmt.Println("Excluded packages:")
	}
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf("Failed to walk in %s due to: %s", path, err)
		}
		if info.IsDir() {
			// Global exclusion of directories
			for _, patt := range strings.Split(excludes, "|") {
				match, err := filepath.Match(patt, filepath.Base(path))
				if err != nil {
					log.Fatalf("Failed to parse pattern: %s\n%v", excludes, err)
				}
				if match {
					if vverbose {
						fmt.Printf("\t%s\n", info.Name())
					}
					return filepath.SkipDir
				}
			}
		} else {
			// Construct the list of files to scan
			// but skip the transformation file if present
			if info.Name() != filepath.Base(tdfPath) {
				files = append(files, path)
			}
		}

		return err
	})

	if vverbose {
		fmt.Println("---")
	}

	if err != nil {
		log.Fatalf("Problem walking to the file or directory:\n %v\n", err)
	}
	return files
}

func shortPath(path string) string {
	wd, err := os.Getwd()
	if err != nil {
		return path
	}

	relPath, err2 := filepath.Rel(wd, path)
	if err2 != nil {
		return path
	}

	return relPath
}

func processFiles(files []string, transformations T) int {
	count := 0
	done := make(chan string, len(files))

	for _, f := range files {

		go func(filePath string) {
			if verbose {
				fmt.Printf("Check file %s\n", shortPath(filePath))
			}

			origDat, data := processFile(filePath, transformations)
			if bytes.Compare(origDat, data) != 0 {

				err := ioutil.WriteFile(filePath, data, 0644)
				if err != nil {
					fmt.Printf("Error writting file %s\n", filePath)
				}

				count++

				if verbose {
					fmt.Printf("Updated file %s\n", shortPath(filePath))
				}

			} else if vverbose {
				fmt.Printf("No update for %s\n", filePath)
			}

			done <- "ok"
		}(f)
	}

	for _ = range files {
		<-done
	}
	if vverbose {
		fmt.Printf("---\n\nChecked %v files\n\n", len(files))
	}
	return count
}

func processFile(filePath string, t T) ([]byte, []byte) {
	var origDat []byte
	var data []byte
	for _, transf := range t.Transformations {
		if checkFileName(filePath, transf) {
			// Initialize the origine data the first time
			if len(origDat) == 0 {
				dat, err := ioutil.ReadFile(filePath)
				if err != nil {
					fmt.Errorf("Error reading file %s\n", filePath)
				}
				data = dat
				origDat = dat
			}

			// If preconditions matche then apply the transformations
			if checkCondition(filePath, data, transf) {
				if vverbose {
					fmt.Printf("Apply tranformation to %s\n", filePath)
				}
				data = applyProcs(data, transf)
			} else {
				if vverbose {
					fmt.Printf("%s doesn't match the preconditions\n", filePath)
				}
			}
		}
	}
	return origDat, data
}
