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

var expectedCount = 6
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

func TestProcessFiles(t *testing.T) {
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
