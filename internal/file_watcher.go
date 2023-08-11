package internal

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"time"
)

var (
	watcher *fsnotify.Watcher
	err     error
)

func waitUntilFind(filename string) (*fsnotify.Event, error) {
	if _, err := os.Stat(filename); err == nil {
		// File already exists when added
		return nil, nil
	}

	for {
		time.Sleep(1 * time.Second)
		if _, err := os.Stat(filename); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}
		break
	}
	// File was created after the add request, so fudge an event
	return &fsnotify.Event{
		Name: filename,
		Op:   fsnotify.Create,
	}, nil
}

func WatchFiles(fileOps chan fsnotify.Event, errs chan error, files ...string) error {
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	for _, filename := range files {
		e, err := waitUntilFind(filename)
		if err != nil {
			return err
		}
		if e != nil {
			fileOps <- *e
		}

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
				log.Printf("new %s event on %s", event.Op, event.Name)
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
