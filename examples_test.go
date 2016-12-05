package watcher_test

import (
	"classes/Servers/Config"
	"log"
	"vendor/github.com/jonathankentstevens/watcher"
)

func ExampleNew() {
	w := watcher.New()
	defer w.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case <-w.Events:
				err := Config.Refresh()
				if err != nil {
					log.Println("Config refresh failed")
				} else {
					log.Println("Config refreshed")
				}
			case err := <-w.Errors:
				log.Println("Config File watcher error:", err.Path, err.Msg)
			}
		}
	}()

	err := w.Add("/path/to/some/file.txt")
	if err != nil {
		log.Fatalln(err)
	}
	<-done
}
