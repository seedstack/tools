// Copyright (c) 2013-2015 by The SeedStack authors. All rights reserved.

// This file is part of SeedStack, An enterprise-oriented full development stack.

// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func walkDir(root string, excludes string, tdfPath string) []string {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			if strings.Contains(excludes, info.Name()) {
				return filepath.SkipDir
			}
		} else {
			if info.Name() != filepath.Base(tdfPath) {
				files = append(files, path)
			}
		}

		return err
	})
	if err != nil {
		log.Fatalf("Problem walking to the file or directory:\n %v\n", err)
	}
	return files
}

func processFiles(files []string, transformations T) int {
	first := true
	count := 0
	done := make(chan string, len(files))

	for _, f := range files {
		go func(filePath string) {
			origDat, data := processFile(filePath, transformations)
			if len(origDat) != len(data) {
				if first && verbose {
					fmt.Println("Apply transformations:")
					first = false
				}
				err := ioutil.WriteFile(filePath, data, 0644)
				count++
				if err != nil {
					fmt.Printf("Error writting file %s\n", filePath)
				}
				if verbose {
					fmt.Printf("%s\n", filePath)
				}
			}

			done <- "ok"
		}(f)
	}

	for _ = range files {
		<-done
	}
	if verbose {
		fmt.Printf("Checked %v files\n", len(files))
	}
	return count
}

func processFile(filePath string, t T) ([]byte, []byte) {
	var origDat []byte
	var data []byte
	for _, transf := range t.Transformations {
		if checkFileName(filePath, transf) {
			if len(origDat) == 0 {
				dat, err := ioutil.ReadFile(filePath)
				if err != nil {
					fmt.Printf("Error reading file %s\n", filePath)
				}
				data = dat
				origDat = dat
			}

			if checkCondition(filePath, data, transf) {
				data = applyProcs(data, transf)
			}
		}
	}
	return origDat, data
}
