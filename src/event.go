package main

import (
	"context"
	"github.com/phiskills/neo4j-client.go"
	"github.com/phisuite/schema.go"
	"log"
)

type eventServer struct {
	schema.UnimplementedEventWriteAPIServer
	store *store
}

func (e *eventServer) Create(_ context.Context, in *schema.Event) (out *schema.Event, err error) {
	var fields []Field
	for _, field := range in.Payload {
		fields = append(fields, Field{
			Name:     field.Name,
			Type:     field.Type.String(),
			Category: field.Category.String(),
		})
	}
	result, err := e.store.CreateData(in.Name, in.Version, fields)
	if err != nil {
		return
	}
	out = e.parse(result, in.Name, in.Payload...)
	log.Printf("[Created] %v", out)
	return
}

func (e *eventServer) Update(_ context.Context, in *schema.Event) (out *schema.Event, err error) {
	var fields []Field
	for _, field := range in.Payload {
		fields = append(fields, Field{
			Name:     field.Name,
			Type:     field.Type.String(),
			Category: field.Category.String(),
		})
	}
	result, err := e.store.UpdateData(in.Name, in.Version, fields)
	if err != nil {
		return
	}
	out = e.parse(result, in.Name, in.Payload...)
	log.Printf("[Updated] %v", out)
	return
}

func (e *eventServer) Activate(_ context.Context, in *schema.Options) (out *schema.Event, err error) {
	status := schema.Status_ACTIVATED.String()
	result, err := e.store.UpdateStatus(status, in.Name, in.Version)
	if err != nil {
		return
	}
	out = e.parse(result, in.Name)
	log.Printf("[Activated] %v", out)
	return
}

func (e *eventServer) Deactivate(_ context.Context, in *schema.Options) (out *schema.Event, err error) {
	status := schema.Status_DEACTIVATED.String()
	result, err := e.store.UpdateStatus(status, in.Name, in.Version)
	if err != nil {
		return
	}
	out = e.parse(result, in.Name)
	log.Printf("[Deactivated] %v", out)
	return
}

func (e *eventServer) parse(event neo4j.Records, name string, payload ...*schema.Field) (out *schema.Event) {
	status := schema.Status_UNACTIVATED
	if s, ok := event[name+".status"]; ok {
		status = schema.Status(schema.Status_value[s.(string)])
	}
	out = &schema.Event{
		Name:    event[name+".name"].(string),
		Version: event[name+".version"].(string),
		Status:  status,
		Payload: []*schema.Field{},
	}
	for _, field := range payload {
		fieldType := event[field.Name+".type"].(string)
		fieldCategory := event[name+"_"+field.Name+".category"].(string)
		out.Payload = append(out.Payload, &schema.Field{
			Name:     event[field.Name+".name"].(string),
			Type:     schema.Field_Type(schema.Field_Type_value[fieldType]),
			Category: schema.Field_Category(schema.Field_Category_value[fieldCategory]),
		})
	}
	return
}
