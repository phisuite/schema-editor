package main

import (
	"context"
	"github.com/phisuite/schema.go"
	"log"
)

type processServer struct {
	schema.UnimplementedProcessAPIServer
}

func (p processServer) Create(_ context.Context,process *schema.Process) (*schema.Process, error) {
	log.Printf("Create: %v", process)
	return process, nil
}

func (p processServer) Update(_ context.Context,process *schema.Process) (*schema.Process, error) {
	log.Printf("Update: %v", process)
	return process, nil
}

func (p processServer) Deactivate(_ context.Context,process *schema.Process) (*schema.Process, error) {
	log.Printf("Deactivate: %v", process)
	process.Status = schema.Status_DEACTIVATED
	return process, nil
}
