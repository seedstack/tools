# SeedStack tools [![Build status](https://travis-ci.org/seedstack/tools.svg?branch=master)](https://travis-ci.org/seedstack/tools)

Tool to perform common tasks on SEED projects.

# Install

```bash
go get github.com/seedstack/tools
```

The following assumes you have Go properly installed.

# Usage

Apply the transformations describe in the `tdf.yml` file to the
current directory. Using the `-t` option of the `fix` subcommand.

```bash
cd $GOPATH/src/github.com/seedstack/tools/test
seed -t tdf.yml fix
```

You can also specify the directory where apply the transformations

```bash
cd $GOPATH/src/github.com/seedstack/tools
seed -t ./test/tdf.yml fix ./test
```

You can also specify the directory where apply the transformation. And
get the transformation file from HTTP as follows.

```bash
cd $GOPATH/src/github.com/seedstack/tools/test
seed -t https://raw.githubusercontent.com/seedstack/tools/master/seed/tdf.yml fix
```

The following assumes you have `$GOPATH/bin` in your `PATH`

# Copyright and license
Code and documentation copyright 2013-2015 The SeedStack authors, released under the MPL 2.0 license.
