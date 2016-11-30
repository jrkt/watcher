//package watcher allows for single & multi-file watching to execute various tasks
package watcher

import (
	"io/ioutil"
	"os"
	"sync"
	"time"
)

const (
	//const identifying an operation describing modification
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

// Event houses the name of the file and the operation performed
type Event struct {
	Name      string
	Operation string
}

// Error contains the error type, path to the file that originated the error,
// the file.Info object and error message
type Error struct {
	error
	Path string
	File File
	Msg  string
}

// File is a wrapper to house the os.FileInfo and the LastSize the file was
type File struct {
	LastSize int64
	Info     os.FileInfo
}

// New returns an empty Watcher pointer
func New() *Watcher {
	w := &Watcher{
		Events: make(chan Event),
		Errors: make(chan Error),
		paths:  make(map[string]File),
		done:   make(chan struct{}),
	}
	return w
}

func (w *Watcher) watch(path string, f File) {
	for {
		info, err := os.Stat(path)
		if err != nil {
			w.Errors <- Error{Path: path, File: f, Msg: err.Error()}
			w.Remove(path)
			return
		} else {
			if f.LastSize != info.Size() {
				w.Events <- Event{Name: path, Operation: OPERATION_MODIFIED}
				f.LastSize = info.Size()
			}
		}
		time.Sleep(1 * time.Second)
	}
}

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
				w.paths[fPath] = File{LastSize: info.Size(), Info: info}
			}
			w.Unlock()
			go w.watch(fPath, w.paths[fPath])
		}
	} else {
		w.Lock()
		_, exists := w.paths[path]
		if !exists {
			w.paths[path] = File{LastSize: info.Size(), Info: info}
		}
		w.Unlock()
		go w.watch(path, w.paths[path])
	}

	return nil
}

func (w *Watcher) Remove(path string) {
	w.Lock()
	defer w.Unlock()
	delete(w.paths, path)
}

func (w *Watcher) Close() {
	close(w.Events)
	close(w.Errors)
	close(w.done)
}
