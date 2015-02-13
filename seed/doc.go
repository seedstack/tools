// Copyright (c) 2013-2015 by The SeedStack authors. All rights reserved.

// This file is part of SeedStack, An enterprise-oriented full development stack.

// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

/*

Perform common tasks on SEED projects.

Install

        go get github.com/seedstack/tools

The following assumes you have Go properly installed.

Usage

Apply the transformations describe in the "tdf.yml" file to the
current directory. Using the "-t" option of the "fix" subcommand.

        cd $GOPATH/src/github.com/seedstack/tools/test
        seed -t tdf.yml fix

You can also specify the directory where apply the transformations

        cd $GOPATH/src/github.com/seedstack/tools
        seed -t ./test/tdf.yml fix ./test

You can also specify the directory where apply the transformation. And
get the transformation file from HTTP as follows.

       cd $GOPATH/src/github.com/seedstack/tools/test
       seed -t https://raw.githubusercontent.com/seedstack/tools/master/seed/tdf.yml fix

The following assumes you have "$GOPATH/bin" in your `PATH`

*/
package main
