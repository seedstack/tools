// Copyright (c) 2013-2015 by The SeedStack authors. All rights reserved.

// This file is part of SeedStack, An enterprise-oriented full development stack.

// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	verbose = false
	vverbose = false
	res := m.Run()
	os.Exit(res)
}

var tdfYml = `exclude: "*.out"
transformations:
 - 
  filter: "*.go|*.yml"
  pre: 
   - AlwaysTrue
  proc:
   - 
    name: Replace
    params:
     - "old"
     - "new"
 - 
  filter: "*.java"
  pre: 
   - AlwaysTrue
  proc:
   - name: DoNothing
`

func TestParseTdf(t *testing.T) {
	tr := parseTdf([]byte(tdfYml), "yml")

	if tr.Exclude != "*.out" {
		t.Error("The file should contains exclude directories.")
	}
	if len(tr.Transformations) != 2 {
		t.Error("The tfl file should contains two transformations.")
	}

	tranf := tr.Transformations[0]

	if tranf.Filter != "*.go|*.yml" {
		t.Error("The first transformation should contains include files.")
	}
	if tranf.Pre[0] != "AlwaysTrue" {
		t.Error("The first transformation should contains a precondition.")
	}
	if len(tranf.Proc) != 1 || tranf.Proc[0].Name != "Replace" || tranf.Proc[0].Params[0] != "old" {
		t.Error("The first transformation should contains a 'replace' procedure.")
	}
}

var tdfToml = `exclude= "*.out"

[[transformations]]
  filter = "*.go|*.yml"
  pre = [ "AlwaysTrue" ]

  [[transformations.proc]]
   name = "Replace"
   params = [ "old", "new" ]

[[transformations]]
  filter = "*.java"
  pre = [ "AlwaysTrue" ]

  [[transformations.proc]]
    name = "DoNothing"
`

func TestParseTdfWithToml(t *testing.T) {
	tr := parseTdf([]byte(tdfToml), "toml")

	if tr.Exclude != "*.out" {
		t.Error("The file should contains exclude directories.")
	}
	if len(tr.Transformations) != 2 {
		t.Error("The tfl file should contains two transformations.")
	}

	tranf := tr.Transformations[0]

	if tranf.Filter != "*.go|*.yml" {
		t.Error("The first transformation should contains include files.")
	}
	if tranf.Pre[0] != "AlwaysTrue" {
		t.Error("The first transformation should contains a precondition.")
	}
	if len(tranf.Proc) != 1 || tranf.Proc[0].Name != "Replace" || tranf.Proc[0].Params[0] != "old" {
		t.Error("The first transformation should contains a 'Replace' procedure.\n%v", tr)
	}
}

func TestGetFormat(t *testing.T) {
	ext, err := getFormat("my/path.yml")
	if err != nil || ext != "yml" {
		t.Errorf("yml was expected but found %s, %v", ext, err)
	}

	ext, err = getFormat("my/path.yaml")
	if err != nil || ext != "yml" {
		t.Errorf("yaml was expected but found %s, %v", ext, err)
	}

	ext, err = getFormat("my/path.toml")
	if err != nil || ext != "toml" {
		t.Errorf("tomml was expected but found %s, %v", ext, err)
	}

	ext, err = getFormat("my/path.TOML")
	if err != nil || ext != "toml" {
		t.Errorf("TOML was expected but found %s, %v", ext, err)
	}

	if _, err := getFormat("my/path.fancy"); err == nil {
		t.Errorf("unsupported format error was expected, but found: %s", err)
	}
}

func TestReadFile(t *testing.T) {
	if bytes := readFile("../test/tdf.yml"); bytes == nil {
		t.Error("ReadFile: Failed to read ./test/conf.yml")
	}

}
