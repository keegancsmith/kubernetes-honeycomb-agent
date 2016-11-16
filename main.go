package main

import (
	"errors"
	"log"
	"os"
	"time"

	libhoney "github.com/honeycombio/libhoney-go"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/rest"
)

func do() error {
	writeKey := os.Getenv("HONEYCOMB_TEAM")
	if writeKey == "" {
		return errors.New("missing required env HONEYCOMB_TEAM")
	}
	dataset := os.Getenv("HONEYCOMB_DATASET")
	if dataset == "" {
		dataset = "kubernetes-events"
	}
	err := libhoney.Init(libhoney.Config{
		WriteKey: writeKey,
		Dataset:  dataset,
	})
	if err != nil {
		return err
	}
	defer libhoney.Close()

	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	for {
		start := time.Now()
		w, err := clientset.Core().Events(api.NamespaceAll).Watch(v1.ListOptions{})
		if err != nil {
			return err
		}
		for {
			watchEvent, ok := <-w.ResultChan()
			if !ok {
				log.Println("events channel closed")
				break
			}
			event := watchEvent.Object.(*api.Event)
			if event.LastTimestamp.Time.Before(start) {
				continue // watch returns old events on startup
			}
			ev := libhoney.NewEvent()
			ev.Timestamp = event.LastTimestamp.Time
			ev.AddField("namespace", event.InvolvedObject.Namespace)
			ev.AddField("kind", event.InvolvedObject.Kind)
			ev.AddField("name", event.InvolvedObject.Name)
			ev.AddField("subobject", event.InvolvedObject.FieldPath)
			ev.AddField("reason", event.Reason)
			ev.AddField("message", event.Message)
			ev.AddField("source_component", event.Source.Component)
			ev.AddField("source_host", event.Source.Host)
			ev.AddField("count", event.Count)
			ev.AddField("type", event.Type)
			// TODO not sure on best format to send up timestamps,
			// using the same format as the event.Timestamp is
			// sent as
			ev.AddField("firstseen", event.FirstTimestamp.Format(time.RFC3339))
			err := ev.Send()
			if err != nil {
				// This should only happen if we screw up
				return err
			}
		}
	}
}

func main() {
	err := do()
	if err != nil {
		log.Fatal(err)
	}
}
