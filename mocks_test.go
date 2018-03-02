package main

import (
	"context"

	"github.com/kkentzo/tagger/indexers"
	"github.com/kkentzo/tagger/watchers"
	"github.com/stretchr/testify/mock"
)

type MockIndexer struct {
	mock.Mock
}

func (indexer *MockIndexer) Create() indexers.Indexable {
	args := indexer.Called()
	return args.Get(0).(indexers.Indexable)
}

func (indexer *MockIndexer) Index(root string, event watchers.Event) {
	indexer.Called(root, event)
}

func (indexer *MockIndexer) CreateWatcher(root string) watchers.Watchable {
	args := indexer.Called(root)
	return args.Get(0).(watchers.Watchable)
}

type MockWatcher struct {
	mock.Mock
}

func CreateMockWatcher() *MockWatcher {
	watcher := &MockWatcher{}
	watcher.On("Watch", mock.AnythingOfType("*context.cancelCtx"))
	watcher.On("Close")
	watcher.On("Events")
	return watcher
}

func (watcher *MockWatcher) Watch(ctx context.Context) {
	watcher.Called(ctx)
}

func (watcher *MockWatcher) Events() chan watchers.Event {
	args := watcher.Called()
	if len(args) > 0 {
		return args.Get(0).(chan watchers.Event)
	} else {
		return make(chan watchers.Event)
	}
}

func (watcher *MockWatcher) Close() {
	watcher.Called()
}
