package internal_test

import (
	"fmt"
	"github.com/dblencowe/k8s-reloader/internal"
	"github.com/fsnotify/fsnotify"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestWatchFiles(t *testing.T) {
	cwd, _ := os.Getwd()
	dir, err := os.MkdirTemp(cwd, "test-")
	if err != nil {
		log.Fatal("test", err)
	}
	defer os.RemoveAll(dir)

	t.Run("Detects changed files", func(t *testing.T) {
		f, err := os.CreateTemp(dir, "watchedfile")
		if err != nil {
			t.Fatal(err)
		}
		defer syscall.Unlink(f.Name())

		fileOps := make(chan fsnotify.Event)
		errChan := make(chan error)
		go func() {
			err := internal.WatchFiles(fileOps, errChan, f.Name())
			assert.NoError(t, err)
		}()

		time.Sleep(3 * time.Second)
		//assert.Equal(t, 0, len(errChan))
		_ = os.WriteFile(f.Name(), []byte("test"), 0644)

		event := <-fileOps
		assert.Equal(t, fsnotify.Write, event.Op)
		assert.Equal(t, f.Name(), event.Name)
	})

	t.Run("Detects created file", func(t *testing.T) {
		fileOps := make(chan fsnotify.Event)
		errChan := make(chan error)
		go func() {
			err := internal.WatchFiles(fileOps, errChan, fmt.Sprintf("%s/%s", dir, "test.txt"))
			assert.NoError(t, err)
		}()

		time.Sleep(3 * time.Second)
		f, err := os.Create(fmt.Sprintf("%s/%s", dir, "test.txt"))
		if err != nil {
			t.Fatal(err)
		}
		defer syscall.Unlink(f.Name())

		event := <-fileOps
		assert.Equal(t, fsnotify.Create, event.Op)
		assert.Equal(t, fmt.Sprintf("%s/%s", dir, "test.txt"), event.Name)
	})

	//t.Run("Detects removed files", func(t *testing.T) {
	//	f, err := os.CreateTemp(dir, "watchedfile")
	//	if err != nil {
	//		t.Fatal(err)
	//	}
	//	fileOps := make(chan fsnotify.Event)
	//	errChan := make(chan error)
	//	go func() {
	//		err := internal.WatchFiles(fileOps, errChan, f.Name())
	//		assert.NoError(t, err)
	//	}()
	//
	//	time.Sleep(3 * time.Second)
	//	syscall.Unlink(f.Name())
	//	time.Sleep(3 * time.Second)
	//	close(fileOps)
	//	assert.Contains(t, fsnotify.Event{
	//		Name: f.Name(),
	//		Op:   fsnotify.Remove,
	//	}, fileOps)
	//})
}
