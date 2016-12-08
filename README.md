[![License](http://img.shields.io/:license-gpl3-blue.svg)](http://www.gnu.org/licenses/gpl-3.0.html)
[![Go Report Card](https://goreportcard.com/badge/github.com/jonathankentstevens/watcher)](https://goreportcard.com/report/github.com/jonathankentstevens/watcher)
[![GoDoc](https://godoc.org/github.com/jonathankentstevens/watcher?status.svg)](https://godoc.org/github.com/jonathankentstevens/watcher)
[![Build Status](https://travis-ci.org/jonathankentstevens/watcher.svg?branch=master)](https://travis-ci.org/jonathankentstevens/watcher)

# watcher
Simple file or directory watcher

# implementation

    go get github.com/jonathankentstevens/watcher

# usage 

```go
package main

import (
	"log"
	"github.com/jonathankentstevens/watcher"
)

func main() {
	w := watcher.NewWatcher()
	defer w.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
                case event := <-w.Events:
                    log.Println("File changed", event)
                case err := <-w.Errors:
                    log.Println("Error Path", err.Path, err.Error())
                }
		}
	}()

	err := w.Add("/path/file1.txt")
	if err != nil {
		log.Fatalln(err)
	}
	err = w.Add("/path/file2.txt")
	if err != nil {
		log.Fatalln(err)
	}
	<-done
}
```

You can also add a directory to the watch list:

```go
err := w.Add("/etc/conf/")
if err != nil {
    log.Fatalln(err)
}
```

# output

```
2016/11/16 21:00:08 File changed {/path/file1.txt MODIFIED}
2016/11/16 21:00:12 File changed {/path/file2.txt MODIFIED}
```
