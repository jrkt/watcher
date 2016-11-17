package watcher

import (
	"sync"
	"os"
)

const (
	OPERATION_MODIFIED = "MODIFIED"
)

// Watcher watches a set of files, delivering events to a channel.
type Watcher struct {
	Events   chan Event
	Errors   chan error
	sync.Mutex
	paths    map[string]File
	done     chan struct{}
}

type Event struct {
	Name string
	Operation string
}

type File struct {
	LastSize int64
	Info os.FileInfo
}

func NewWatcher() *Watcher {
	w := &Watcher{
		Events:   make(chan Event),
		Errors:   make(chan error),
		paths:    make(map[string]File),
		done:     make(chan struct{}),
	}
	return w
}

func (w *Watcher) watch(path string, f File) {
	for {
		info, err := os.Stat(path)
		if err != nil {
			w.Errors <- err
		} else {
			if f.LastSize != info.Size() {
				w.Events <- Event{Name: path, Operation: OPERATION_MODIFIED}
				f.LastSize = info.Size()
			}
		}
	}
}

func (w *Watcher) Add(path string) error {
	info, err := os.Stat(path);
	if err != nil {
		return err
	}
	w.Lock()
	_, exists := w.paths[path]
	if !exists {
		w.paths[path] = File{LastSize: info.Size(), Info: info}
	}
	w.Unlock()

	go w.watch(path, w.paths[path])

	return nil
}

func (w *Watcher) Close() {
	close(w.Events)
	close(w.Errors)
	close(w.done)
}