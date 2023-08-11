package internal

import (
	"github.com/fsnotify/fsnotify"
	"log"
)

var (
	watcher *fsnotify.Watcher
	err     error
)

func WatchFiles(fileOps chan fsnotify.Event, errs chan error, files ...string) error {
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	for _, filename := range files {
		err = watcher.Add(filename)
		if err != nil {
			return err
		}
		log.Printf("added %s to watched files list", filename)
	}

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Printf("new file event: %+v", event)
				fileOps <- event
			case err := <-watcher.Errors:
				errs <- err
			}
		}
	}()

	return nil
}

func Shutdown() error {
	if watcher != nil {
		return watcher.Close()
	}
	return nil
}
