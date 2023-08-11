package internal_test

import (
	"github.com/dblencowe/k8s-reloader/internal"
	"github.com/fsnotify/fsnotify"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"syscall"
	"testing"
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
		err = internal.WatchFiles(fileOps, errChan, f.Name())
		assert.NoError(t, err)
		_ = os.WriteFile(f.Name(), []byte("test"), 0644)

		event := <-fileOps
		assert.Equal(t, fsnotify.Write, event.Op)
		assert.Equal(t, f.Name(), event.Name)
	})

	//t.Run("Detects created file", func(t *testing.T) {
	//	fileOps := make(chan fsnotify.Event)
	//	errChan := make(chan error)
	//	err := internal.WatchFiles(fileOps, errChan, "watchedfile")
	//	assert.NoError(t, err)
	//	f, err := os.CreateTemp(dir, "watchedfile")
	//	if err != nil {
	//		t.Fatal(err)
	//	}
	//	defer syscall.Unlink(f.Name())
	//
	//	event := <-fileOps
	//	assert.Equal(t, fsnotify.Create, event.Op)
	//	assert.Equal(t, "watchedfile", event.Name)
	//})

	//t.Run("Detects removed files", func(t *testing.T) {
	//	f, err := os.CreateTemp(dir, "watchedfile")
	//	if err != nil {
	//		t.Fatal(err)
	//	}
	//	fileOps := make(chan fsnotify.Event)
	//	errChan := make(chan error)
	//	err = internal.WatchFiles(fileOps, errChan, f.Name())
	//	assert.NoError(t, err)
	//	syscall.Unlink(f.Name())
	//
	//	event := <-fileOps
	//	removed := false
	//	assert.Equal(t, fsnotify.Remove, event.Op)
	//	assert.Equal(t, f.Name(), event.Name)
	//})
}
