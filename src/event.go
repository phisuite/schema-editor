package main

import (
	"context"
	"github.com/phisuite/schema.go"
	"log"
)

type eventServer struct {
	schema.UnimplementedEventWriteAPIServer
}

func (e eventServer) Create(_ context.Context, event *schema.Event) (*schema.Event, error) {
	log.Printf("Create: %v", event)
	return event, nil
}

func (e eventServer) Update(_ context.Context, event *schema.Event) (*schema.Event, error) {
	log.Printf("Update: %v", event)
	return event, nil
}

func (e eventServer) Activate(_ context.Context, options *schema.Options) (*schema.Event, error) {
	event := &schema.Event{Name:options.Name, Version:options.Version}
	log.Printf("Activate: %v", event)
	event.Status = schema.Status_ACTIVATED
	return event, nil
}

func (e eventServer) Deactivate(_ context.Context, options *schema.Options) (*schema.Event, error) {
	event := &schema.Event{Name:options.Name, Version:options.Version}
	log.Printf("Deactivate: %v", event)
	event.Status = schema.Status_DEACTIVATED
	return event, nil
}
