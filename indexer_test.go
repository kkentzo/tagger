package main

import (
	"github.com/stretchr/testify/mock"
)

type MockIndexer struct {
	mock.Mock
}

func (indexer *MockIndexer) Index(root string) {
	indexer.Called(root)
}

func (indexer *MockIndexer) CreateWatcher(root string) Watcher {
	args := indexer.Called(root)
	return args.Get(0).(Watcher)
}
