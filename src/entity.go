package main

import (
	"context"
	"github.com/phisuite/schema.go"
	"log"
)

type entityServer struct {
	schema.UnimplementedEntityWriteAPIServer
}

func (e entityServer) Create(_ context.Context, entity *schema.Entity) (*schema.Entity, error) {
	log.Printf("Create: %v", entity)
	return entity, nil
}

func (e entityServer) Update(_ context.Context, entity *schema.Entity) (*schema.Entity, error) {
	log.Printf("Update: %v", entity)
	return entity, nil
}

func (e entityServer) Activate(_ context.Context, options *schema.Options) (*schema.Entity, error) {
	entity := &schema.Entity{Name:options.Name, Version:options.Version}
	log.Printf("Activate: %v", entity)
	entity.Status = schema.Status_ACTIVATED
	return entity, nil
}

func (e entityServer) Deactivate(_ context.Context, options *schema.Options) (*schema.Entity, error) {
	entity := &schema.Entity{Name:options.Name, Version:options.Version}
	log.Printf("Deactivate: %v", entity)
	entity.Status = schema.Status_DEACTIVATED
	return entity, nil
}
