package main

type Indexer interface {
	Index(string)
	CreateWatcher(string) *Watcher
}

type IndexerType string

const (
	Generic IndexerType = "generic"
	Rvm                 = "rvm"
)
