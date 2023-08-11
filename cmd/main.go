package main

import (
	"context"
	"github.com/dblencowe/k8s-reloader/internal"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
)

func main() {
	files := os.Args[1:]

	deploymentName := requiredEnv("DEPLOYMENT_NAME")
	namespace := requiredEnv("DEPLOYMENT_NAMESPACE")

	_, isDev := os.LookupEnv("IS_DEV")
	k8Client, err := internal.MakeK8Client(isDev)
	if err != nil {
		log.Fatal(err)
	}

	fileOps := make(chan fsnotify.Event)
	errChan := make(chan error)
	err = internal.WatchFiles(fileOps, errChan, files...)
	if err != nil {
		log.Fatal(err)
	}
	defer internal.Shutdown()

	for {
		select {
		case event := <-fileOps:
			log.Printf("File Op: %+v", event)
			err = k8Client.RestartDeployment(context.TODO(), namespace, deploymentName)
			if err != nil {
				log.Fatal(err)
			}
		case err := <-errChan:
			log.Fatal(err)
		}
	}
}

func requiredEnv(key string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	log.Fatalf("Required EnvVar %s not set", key)
	return ""
}
