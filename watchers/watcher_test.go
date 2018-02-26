package watchers

import (
	"context"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_NewWatcher(t *testing.T) {
	watcher := NewWatcher("foo", []string{"excl"}, "TAGS", 2*time.Second)
	defer watcher.Close()
	assert.Equal(t, "foo", watcher.Root)
	assert.Equal(t, 2*time.Second, watcher.MaxPeriod)
	assert.IsType(t, &FsWatcher{}, watcher.fsWatcher)
	assert.IsType(t, make(chan Event), watcher.events)
}

func Test_Watcher_Events_ReturnsTheChannel(t *testing.T) {
	watcher := NewWatcher("foo", []string{"excl"}, "TAGS", 2*time.Second)
	defer watcher.Close()
	go func(w *Watcher) { watcher.Events() <- Event{} }(watcher)
	assert.Equal(t, Event{}, <-watcher.Events())
}

func Test_Watcher_Close_ClosesTheChannels(t *testing.T) {
	watcher := NewWatcher("foo", []string{"excl"}, "TAGS", 2*time.Second)
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

	watcher := NewWatcher("foo", []string{}, "TAGS", 2*time.Second)
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

func Test_Watcher_Watch_ShouldSetSpecialFile_OnFsNotify_Event(t *testing.T) {
	fsWatcher := &MockFsWatcher{}
	events := make(chan fsnotify.Event)
	fsWatcher.On("Events").Return(events)
	fsWatcher.On("Errors").Return(make(chan error))
	fsWatcher.On("Add", "foo").Return(nil)

	watcher := NewWatcher("foo", []string{}, "TAGS", 10*time.Millisecond)
	watcher.fsWatcher = fsWatcher
	watcher.SpecialFile = "the_special_file"

	e := fsnotify.Event{
		Name: "the_special_file",
		Op:   fsnotify.Create,
	}
	fsWatcher.On("Handle", mock.AnythingOfType("fsnotify.Event")).Return(true)

	go watcher.Watch(context.Background())
	// fire the filesystem event
	events <- e
	// expectation
	event := <-watcher.events
	assert.True(t, event.IsSpecial)

}

func Test_Watcher_Watch_ShouldReindex_WhenTickerTicks(t *testing.T) {
	fsWatcher := &MockFsWatcher{}
	events := make(chan fsnotify.Event)
	fsWatcher.On("Events").Return(events)
	fsWatcher.On("Errors").Return(make(chan error))
	fsWatcher.On("Add", "foo").Return(nil)

	watcher := NewWatcher("foo", []string{}, "TAGS", 10*time.Millisecond)
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
	assert.Equal(t, Event{}, <-watcher.events)
}
