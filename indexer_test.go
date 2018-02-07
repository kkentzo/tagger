package main

import (
	"testing"

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

func Test_Indexer_Index_ShouldTriggerCommand(t *testing.T) {
	t.Skip("Pending")
}

func Test_Indexer_CreateWatcher_ShouldReturnAWatcher(t *testing.T) {
	t.Skip("Pending")
}

func Test_Indexer_GetArguments_WhenIndexerIsGeneric(t *testing.T) {
	t.Skip("Pending")
}

func Test_Indexer_GetArguments_WhenIndexerIsRvm(t *testing.T) {
	t.Skip("Pending")
}
