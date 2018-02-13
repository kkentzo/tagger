package main

import (
	"context"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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

func (watcher *MockWatcher) Events() chan struct{} {
	args := watcher.Called()
	if len(args) > 0 {
		return args.Get(0).(chan struct{})
	} else {
		return make(chan struct{})
	}
}

func (watcher *MockWatcher) Close() {
	watcher.Called()
}

func Test_NewWatcher(t *testing.T) {
	watcher := NewWatcher("foo", []string{"excl"}, 2*time.Second)
	defer watcher.Close()
	assert.Equal(t, "foo", watcher.Root)
	assert.Equal(t, 2*time.Second, watcher.MaxFrequency)
	assert.IsType(t, &FsWatcher{}, watcher.fsWatcher)
	assert.IsType(t, make(chan struct{}), watcher.events)
}

func Test_Watcher_Events_ReturnsTheChannel(t *testing.T) {
	watcher := NewWatcher("foo", []string{"excl"}, 2*time.Second)
	defer watcher.Close()
	var s struct{}
	go func(w *Watcher) { watcher.Events() <- s }(watcher)
	assert.Equal(t, s, <-watcher.Events())
}

func Test_Watcher_Close_ClosesTheChannels(t *testing.T) {
	watcher := NewWatcher("foo", []string{"excl"}, 2*time.Second)
	watcher.Close()
	_, open := <-watcher.Events()
	assert.False(t, open)
	_, open = <-watcher.fsWatcher.Events()
	assert.False(t, open)
}

func Test_Watcher_Watch_ShouldCallHandler_OnFsNotify_Event(t *testing.T) {
	fsWatcher := &MockFsWatcher{}
	events := make(chan fsnotify.Event)
	fsWatcher.On("Events").Return(events)
	fsWatcher.On("Errors").Return(make(chan error))
	fsWatcher.On("Add", "foo").Return(nil)

	watcher := NewWatcher("foo", []string{}, 2*time.Second)
	watcher.fsWatcher = fsWatcher

	e := fsnotify.Event{
		Name: "foo",
		Op:   fsnotify.Create,
	}
	fired := make(chan fsnotify.Event)
	fsWatcher.
		On("Handle", mock.AnythingOfType("fsnotify.Event")).
		Return(true).
		Run(func(args mock.Arguments) {
			fired <- e
		})

	go watcher.Watch(context.Background())
	// fire the filesystem event
	events <- e
	// expectation
	assert.Equal(t, e, <-fired)
}

func Test_Watcher_Watch_ShouldReindex_WhenTickerTicks(t *testing.T) {
	fsWatcher := &MockFsWatcher{}
	events := make(chan fsnotify.Event)
	fsWatcher.On("Events").Return(events)
	fsWatcher.On("Errors").Return(make(chan error))
	fsWatcher.On("Add", "foo").Return(nil)

	watcher := NewWatcher("foo", []string{}, 10*time.Millisecond)
	watcher.fsWatcher = fsWatcher

	e := fsnotify.Event{
		Name: "foo",
		Op:   fsnotify.Create,
	}
	fsWatcher.On("Handle", mock.AnythingOfType("fsnotify.Event")).Return(true)

	go watcher.Watch(context.Background())
	// fire the filesystem event
	events <- e
	// expectation
	var msg struct{}
	assert.Equal(t, msg, <-watcher.events)
}
