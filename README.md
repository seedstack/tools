# SeedStack tools [![Build status](https://travis-ci.org/seedstack/tools.svg?branch=master)](https://travis-ci.org/seedstack/tools)

Tool to perform common task on SEED projects.

# Install

```bash
go get github.com/seedstack/tools
```

The following assumes you have Go properly installed.

# Usage

Apply the transformation in the `tansformation.yml` file to the
current directory.

```bash
cd $GOPATH/src/github.com/seedstack/tools/test
seed -t tdf.yml
```

You can also specify the directory where apply the transformation. And
get the transformation file from HTTP as follows.

```bash
cd $GOPATH/src/github.com/seedstack/tools
seed -t https://raw.githubusercontent.com/seedstack/tools/master/seed/tdf.yml ./test
```

The following assumes you have `$GOPATH/bin` in your `PATH`

# Copyright and license
Code and documentation copyright 2013-2015 The SeedStack authors, released under the MPL 2.0 license.
