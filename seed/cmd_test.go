// Copyright (c) 2013-2015 by The SeedStack authors. All rights reserved.

// This file is part of SeedStack, An enterprise-oriented full development stack.

// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"testing"
)

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
	tr := parseTdf([]byte(tdfYml))

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
