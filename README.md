[![Travis Build
Status](https://travis-ci.org/kkentzo/tagger.svg?branch=master)](https://travis-ci.org/kkentzo/tagger)
[![Go Report
Card](https://goreportcard.com/badge/github.com/kkentzo/tagger)](https://goreportcard.com/report/github.com/kkentzo/tagger)

# tagger

`tagger` is a program that monitors multiple separate project
directory hierarchies and triggers a particular command when a file
change is detected within these directory structures.

The main use case (and the one for which `tagger` was developed) is to
facilitate the automatic and dynamic indexing of multiple software
projects using an external program such as `ctags`. While the use of
`ctags` is quite popular among `emacs` and `vi` users especially for
dynamically-typed languages, there exist no tools for automating the
re-indexing of multiple projects as their respectives files are added,
removed and modified. `tagger` is developed in order to fill this gap.

# Features

`tagger` has the following features:

* recursive monitoring of multiple project directories (using
  `fsnotify`)
* configuration over the program and arguments to run when a file
  change is detected (default is `ctags -R -e`)
* exclusion filters for ignoring project directories
* some support for secondary project directories (libraries) that are
  located outside the project directory tree (only ruby's
  `rvm`/`bundler` gemset paths supported at the moment)
* throttling of reindexing events; especially useful for actively
  developed projects
* a `yaml` configuration file for statically specifying which projects
  to monitor
* an http interface for adding/removing/listing projects dynamically at runtime

# Known Issues

On MacOS there exists a [known
issue](https://github.com/fsnotify/fsnotify/issues/129) with
filesystem monitoring using kqueue through the fsnotify library that
produces a "too many open files" error. The solution is probably some
adjustment of the relevant OS limit.

# Installation

Linux binaries are available from the project's [release
page](https://github.com/kkentzo/tagger/releases).

# Usage

`tagger -h` will provide usage details for the program. A sample
configuration yaml is also [provided for reference](demo.yml) for
specifying the list of projects that `tagger` will start monitoring as
well as indexer-specific details.

# Development and Tests

First of all, make sure that you have a [working go
installation](https://golang.org/doc/install) (this includes a valid
`$GOPATH`). `tagger` can then be installed using:

``` bash
$ go get -u github.com/kkentzo/tagger
```

The `dep` tool is also necessary for installing the program's
dependencies:

``` bash
$ go get -u github.com/golang/dep/cmd/dep
$ cd $GOPATH/src/github.com/kkentzo/tagger
$ dep ensure
```

The standard go tooling can be used to build the project and run the
tests:

``` bash
$ go build
$ go test ./... -v
```

# Wish list

Some features that may be added in the project at some point:

* generalize a project's ancillary (secondary) indexer beyond
  ruby/rvm(for libraries that are located outside the project's
  directory tree)
* implement adaptive indexing by dynamically measuring indexing time
  per project and adjusting throttling accordingly

# Contributing

Bug reports and pull requests are welcome on GitHub at the [project's
page](https://github.com/kkentzo/tagger). This project is intended to
be a safe, welcoming space for collaboration, and contributors are
expected to adhere to the [Contributor
Covenant](http://contributor-covenant.org) code of conduct.

# License

`tagger` is available as open source under the terms of the [MIT
License](http://opensource.org/licenses/MIT).
