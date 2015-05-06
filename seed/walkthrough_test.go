// Copyright (c) 2013-2015 by The SeedStack authors. All rights reserved.

// This file is part of SeedStack, An enterprise-oriented full development stack.

// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"path/filepath"
	"testing"
)

const expectedCount = 6

var expectedFile = filepath.FromSlash("../test/dir1/file21")

func TestWalkDir(t *testing.T) {
	files := walkDir("../test", "", "../test/tdf.yml")
	if len(files) != expectedCount {
		t.Errorf("WalkDir expect %v files but found %v", expectedCount, len(files))
	}

	if files[0] != expectedFile {
		t.Errorf("WalkDir expect %v but found %v", expectedFile, files[0])
	}

	files = walkDir("../test", "test", "../test/tdf.yml")
	if len(files) != 0 {
		t.Errorf("WalkDir expect %v files but found %v", 0, len(files))
	}
}

func TestShortPath(t *testing.T) {
	originalPath := filepath.Join("..","test","dir1","file21")
	sp := shortPath(expectedFile)

	if sp != originalPath {
		t.Errorf("shorpath: %s was expected but found %s", originalPath, sp)
	}
}

func TestProcessFiles(t *testing.T) {
	p := []Procedure{Procedure{Name: "Insert", Params: []string{"foo"}}}
	tt := Transformation{Filter: "*file1", Proc: p}
	tf := Transformation{Filter: "*.go", Proc: p}
	filesToCheck := []string{"../test/file1", "../test/file1", "../test/file2"}
	expectedCount := 2
	
	modifiedFiles := processFiles(filesToCheck, T{Transformations: []Transformation{tt, tf}})

	if modifiedFiles != expectedCount {
		t.Errorf("processFiles: %v files should be processed but found %v", expectedCount, modifiedFiles)
	}

	modifiedFiles = processFiles(filesToCheck, T{Transformations: []Transformation{}})

	if modifiedFiles != 0 {
		t.Errorf("processFiles: no files should be processed but found %v", expectedCount, modifiedFiles)
	}

	// Cleanup
	r := []Procedure{Procedure{Name: "RemoveAtEnd", Params: []string{"foo"}}}
	cleanup := Transformation{Filter: "*file1", Proc: r}
	filesToClean := []string{"../test/file1", "../test/file1"}
	processFiles(filesToClean, T{Transformations: []Transformation{cleanup}})
}

func TestProcessFile(t *testing.T) {
	p := []Procedure{Procedure{Name: "Insert", Params: []string{"foo"}}}
	tt := Transformation{Filter: "*file1", Proc: p}
	tf := Transformation{Filter: "*.go", Proc: p}

	orig, dat := processFile("../test/file1", T{Transformations: []Transformation{tt}})
	if string(orig) == string(dat) {
		t.Error("file1 should be processed.")
	}

	orig, dat = processFile("../test/file1", T{Transformations: []Transformation{tf}})
	if string(orig) != string(dat) {
		t.Error("file1 should not be processed.")
	}
}
