SeedStack tools
====
[![Build status](https://travis-ci.org/seedstack/tools.svg?branch=master)](https://travis-ci.org/seedstack/tools) [![Coverage Status](https://coveralls.io/repos/seedstack/tools/badge.svg?branch=master)](https://coveralls.io/r/seedstack/tools?branch=master)

Tool to perform common tasks on SEED projects. The main provided
command is `seed fix`. Type `seed help fix` to see a more detailed
documentation.

# Download

Download the [latest version](https://github.com/seedstack/tools/releases).

# Install from source

The following assumes you have Go properly installed and that you have
`$GOPATH/bin` in your `PATH`.

```bash
go get github.com/seedstack/tools/seed
seed
```

# Usage

Apply the transformations described in the `transform.toml` file to
the current directory. Use the `-t` option of the `fix` subcommand.

```bash
seed -t transform.toml fix
```

You can specify the directory where to apply the transformations:

```bash
seed -t ./test/transform.toml fix ./myproject
```

You can also retrieve the transformation file from HTTP:

```bash
seed -t https://raw.githubusercontent.com/seedstack/tools/master/seed/tdf.yml fix
```

# Copyright and license

Code and documentation copyright 2013-2015 The SeedStack authors,
released under the MPL 2.0 license.
