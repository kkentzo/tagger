package main

type Indexer struct {
	command string // the tags binary
	args    string // args to ctags binary
}

func (indexer *Indexer) Index(root string) {
	// TODO: execute command in the context of root path
}
