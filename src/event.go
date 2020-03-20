package main

import (
	"context"
	"github.com/phisuite/schema.go"
	"log"
)

type eventServer struct {
	schema.UnimplementedEventAPIServer
}

func (e eventServer) Create(_ context.Context, event *schema.Event) (*schema.Event, error) {
	log.Printf("Create: %v", event)
	return event, nil
}

func (e eventServer) Update(_ context.Context, event *schema.Event) (*schema.Event, error) {
	log.Printf("Update: %v", event)
	return event, nil
}

func (e eventServer) Activate(_ context.Context, event *schema.Event) (*schema.Event, error) {
	log.Printf("Activate: %v", event)
	event.Status = schema.Status_ACTIVATED
	return event, nil
}

func (e eventServer) Deactivate(_ context.Context, event *schema.Event) (*schema.Event, error) {
	log.Printf("Deactivate: %v", event)
	event.Status = schema.Status_DEACTIVATED
	return event, nil
}
