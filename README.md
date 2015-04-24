# SeedStack tools [![Build status](https://travis-ci.org/seedstack/tools.svg?branch=master)](https://travis-ci.org/seedstack/tools) [![Stories in Ready](https://badge.waffle.io/seedstack/tools.png?label=ready&title=Ready)](https://waffle.io/seedstack/tools)

Tool to perform common tasks on SEED projects.

# Install

```bash
go get github.com/seedstack/tools/seed
```

The following assumes you have Go properly installed and that you have `$GOPATH/bin` in your `PATH`.

# Usage

Apply the transformations described in the `tdf.yml` file to the
current directory. Use the `-t` option of the `fix` subcommand.

```bash
cd $GOPATH/src/github.com/seedstack/tools/test
seed -t tdf.yml fix
```

You can specify the directory where to apply the transformations:

```bash
cd $GOPATH/src/github.com/seedstack/tools
seed -t ./test/tdf.yml fix ./test
```

You can retrieve the transformation descriptor from HTTP:

```bash
cd $GOPATH/src/github.com/seedstack/tools/test
seed -t https://raw.githubusercontent.com/seedstack/tools/master/seed/tdf.yml fix
```

# Copyright and license
Code and documentation copyright 2013-2015 The SeedStack authors, released under the MPL 2.0 license.
