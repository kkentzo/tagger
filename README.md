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
re-indexing of multiple projects. `tagger` is developed in order to
fill this gap in tooling so as to make developer lives more productive
and happy.

# Features

`tagger` has the following features:

* recursive monitoring of multiple project directories (using
  `fsnotify`)
* configuration over the program and arguments to run when a file
  change is detected
* exclusion filters for ignoring certain directories
* support for secondary project directories (libraries) that are
  located outside the project directory tree
  * only ruby's `rvm` gemset paths supported at the moment
* throttling of reindexing events; especially useful for actively developed projects
* a `yaml` configuration file for statically specifying which projects
  to monitor
* an http interface for adding/removing/listing projects dynamically at runtime

# Known Issues

On MacOS there exists a [known
issue](https://github.com/fsnotify/fsnotify/issues/129) with
filesystem monitoring using kqueue through the fsnotify library that
produces a "too many open files" error.

# Installation

Linux binaries are available from the [project release page](https://github.com/kkentzo/tagger/releases).

# Usage

# Development and Tests

# Coming Up
