// Copyright (c) 2013-2015 by The SeedStack authors. All rights reserved.

// This file is part of SeedStack, An enterprise-oriented full development stack.

// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"testing"
	"path/filepath"
)

var expectedCount = 6
var expectedFile = filepath.FromSlash("../test/dir1/file21")

func TestWalkthroughDir(t *testing.T) {
	files := walkthroughDir("../test")
	if len(files) != expectedCount {
		t.Errorf("Walkthrough expect %v files but found %v", expectedCount, len(files))
	}
	
	if files[0] != expectedFile {
		t.Errorf("Walkthrough expect %v but found %v", expectedFile, files[0])
	}
}
