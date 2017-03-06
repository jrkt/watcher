// Package watcher is a simple file and directory watcher
package watcher

import (
	"io/ioutil"
	"os"
	"sync"
	"time"
)

const (
	// OPERATION_MODIFIED is the constant defining a modify action
	OPERATION_MODIFIED = "MODIFIED"
)

// Watcher watches a set of files, delivering events to a channel.
type Watcher struct {
	Events chan Event
	Errors chan Error
	sync.Mutex
	paths map[string]File
	done  chan struct{}
}

// Event defines the filename & specific operation performed
type Event struct {
	Name      string
	Env       string
	Operation string
}

// Error holds the error as well as the file path, the os.Info of the file, and the error message
type Error struct {
	error
	Path string
	File File
	Msg  string
}

// File contains each file's last modification time and it's respective os.Info
type File struct {
	LastModTime time.Time
	Info        os.FileInfo
}

// New returns a pointer to a new Watcher instance with everything initialized
func New() *Watcher {
	w := &Watcher{
		Events: make(chan Event),
		Errors: make(chan Error),
		paths:  make(map[string]File),
		done:   make(chan struct{}),
	}
	return w
}

// Add adds a file/directory path to the watch list
func (w *Watcher) Add(path string) error {
	info, err := os.Stat(path)
	if err != nil {

		return err
	}

	if info.IsDir() {
		files, err := ioutil.ReadDir(path)
		if err != nil {
			return err
		}
		for _, file := range files {
			fPath := path + file.Name()
			info, err := os.Stat(path + file.Name())
			if err != nil {
				println("Error adding file to watch list", path+file.Name())
				continue
			}
			w.Lock()
			_, exists := w.paths[fPath]
			if !exists {
				w.paths[fPath] = File{LastModTime: info.ModTime(), Info: info}
			}
			w.Unlock()
			go w.watchFile(fPath, w.paths[fPath])
		}
	} else {
		w.Lock()
		_, exists := w.paths[path]
		if !exists {
			w.paths[path] = File{LastModTime: info.ModTime(), Info: info}
		}
		w.Unlock()
		go w.watchFile(path, w.paths[path])
	}

	return nil
}

// watch serves as a goroutine for each path added to the watch list
func (w *Watcher) watchFile(path string, f File) {
	for {
		info, err := os.Stat(path)
		if err != nil {
			w.Errors <- Error{Path: path, File: f, Msg: err.Error()}
			w.Remove(path)
			return
		}
		modTime := info.ModTime()
		if f.LastModTime != modTime {
			w.Events <- Event{Name: path, Operation: OPERATION_MODIFIED}
			f.LastModTime = modTime
		}
		time.Sleep(1 * time.Second)
	}
}

// Remove removes a path from the watch list
func (w *Watcher) Remove(path string) {
	w.Lock()
	delete(w.paths, path)
	w.Unlock()
}

// Close closes all channels
func (w *Watcher) Close() {
	close(w.Events)
	close(w.Errors)
	close(w.done)
}
