package main

import (
	"context"
	"github.com/phisuite/schema.go"
	"log"
)

type entityServer struct {
	schema.UnimplementedEntityAPIServer
}

func (e entityServer) Create(_ context.Context, entity *schema.Entity) (*schema.Entity, error) {
	log.Printf("Create: %v", entity)
	return entity, nil
}

func (e entityServer) Update(_ context.Context, entity *schema.Entity) (*schema.Entity, error) {
	log.Printf("Update: %v", entity)
	return entity, nil
}

func (e entityServer) Activate(_ context.Context, entity *schema.Entity) (*schema.Entity, error) {
	log.Printf("Activate: %v", entity)
	entity.Status = schema.Status_ACTIVATED
	return entity, nil
}

func (e entityServer) Deactivate(_ context.Context, entity *schema.Entity) (*schema.Entity, error) {
	log.Printf("Deactivate: %v", entity)
	entity.Status = schema.Status_DEACTIVATED
	return entity, nil
}
