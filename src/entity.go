package main

import (
	"context"
	"github.com/phiskills/neo4j-client.go"
	"github.com/phisuite/schema.go"
	"log"
)

type entityServer struct {
	schema.UnimplementedEntityWriteAPIServer
	store *store
}

func (e *entityServer) Create(_ context.Context, in *schema.Entity) (out *schema.Entity, err error) {
	var fields []Field
	for _, field := range in.Data {
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
	out = e.parse(result, in.Name, in.Data...)
	log.Printf("[Created] %v", out)
	return
}

func (e *entityServer) Update(_ context.Context, in *schema.Entity) (out *schema.Entity, err error) {
	var fields []Field
	for _, field := range in.Data {
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
	out = e.parse(result, in.Name, in.Data...)
	log.Printf("[Updated] %v", out)
	return
}

func (e *entityServer) Activate(_ context.Context, in *schema.Options) (out *schema.Entity, err error) {
	status := schema.Status_ACTIVATED.String()
	result, err := e.store.UpdateStatus(status, in.Name, in.Version)
	if err != nil {
		return
	}
	out = e.parse(result, in.Name)
	log.Printf("[Activated] %v", out)
	return
}

func (e *entityServer) Deactivate(_ context.Context, in *schema.Options) (out *schema.Entity, err error) {
	status := schema.Status_DEACTIVATED.String()
	result, err := e.store.UpdateStatus(status, in.Name, in.Version)
	if err != nil {
		return
	}
	out = e.parse(result, in.Name)
	log.Printf("[Deactivated] %v", out)
	return
}

func (e *entityServer) parse(result neo4j.Records, name string, payload ...*schema.Field) (out *schema.Entity) {
	status := schema.Status_UNACTIVATED
	if s, ok := result[name+".status"]; ok {
		status = schema.Status(schema.Status_value[s.(string)])
	}
	out = &schema.Entity{
		Name:    result[name+".name"].(string),
		Version: result[name+".version"].(string),
		Status:  status,
		Data:    []*schema.Field{},
	}
	for _, field := range payload {
		fieldType := result[field.Name+".type"].(string)
		fieldCategory := result[name+"_"+field.Name+".category"].(string)
		out.Data = append(out.Data, &schema.Field{
			Name:     result[field.Name+".name"].(string),
			Type:     schema.Field_Type(schema.Field_Type_value[fieldType]),
			Category: schema.Field_Category(schema.Field_Category_value[fieldCategory]),
		})
	}
	return
}
