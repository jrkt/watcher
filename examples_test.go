package watcher_test

import (
	"log"

	"github.com/jrkt/watcher"
)

func ExampleNew() {
	w := watcher.New()
	defer w.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case ev := <-w.Events:
				log.Println("File changed:", ev.Name)
			case err := <-w.Errors:
				log.Println("Watcher error:", err.Path, err.Msg)
			}
		}
	}()

	err := w.Add("/path/to/some/file.txt")
	if err != nil {
		log.Fatalln(err)
	}
	<-done
}
