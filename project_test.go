package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Project_DefaultProject(t *testing.T) {
	indexer := &MockIndexer{}
	watcher := &MockWatcher{}

	project := DefaultProject(indexer, watcher)

	assert.Equal(t, ".", project.Path)
	assert.Equal(t, indexer, project.Indexer)
	assert.Equal(t, watcher, project.Watcher)
}

func Test_Project_Monitor_WillIndexProject_OnWatcherEvent(t *testing.T) {
	indexer := &MockIndexer{}
	watcher := &MockWatcher{}

	var event struct{}
	events := make(chan struct{})
	go func(e chan struct{}) { events <- event }(events)

	watcher.On("Events").Return(events)
	watcher.On("Watch", mock.AnythingOfType("*context.cancelCtx"))

	indexed := make(chan struct{})
	indexer.On("Index", ".").Run(func(args mock.Arguments) { indexed <- event })

	project := DefaultProject(indexer, watcher)
	go project.Monitor(context.Background())
	<-indexed
}

func Test_Project_Monitor_WillCloseWatcher_OnContextCancellation(t *testing.T) {
	indexer := &MockIndexer{}
	watcher := &MockWatcher{}
	indexer.On("Index", ".").Return()

	events := make(chan struct{})
	watcher.On("Watch", mock.AnythingOfType("*context.cancelCtx"))
	watcher.On("Events").Return(events)
	closed := make(chan struct{})
	watcher.On("Close").Run(func(args mock.Arguments) { close(closed) })

	project := DefaultProject(indexer, watcher)
	ctx, cancel := context.WithCancel(context.Background())
	go cancel()
	project.Monitor(ctx)
	<-closed
}
