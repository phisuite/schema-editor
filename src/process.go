package main

import (
	"context"
	"fmt"
	"github.com/phiskills/neo4j-client.go"
	"github.com/phisuite/schema.go"
	"log"
)

type processServer struct {
	schema.UnimplementedProcessWriteAPIServer
	store *store
}

type processCategory string
const (
	processInput  processCategory = "Input"
	processOutput processCategory = "Output"
	processError  processCategory = "Error"
)

func (p *processServer) Create(_ context.Context, in *schema.Process) (out *schema.Process, err error) {
	definition := p.extract(in.Definition.Input, in.Definition.Output, in.Definition.Error)
	result, err := p.store.CreateProcess(in.Name, in.Version, definition)
	if err != nil {
		return
	}
	out = p.parse(result, in.Name)
	log.Printf("[Created] %v", out)
	return
}

func (p *processServer) Update(_ context.Context, in *schema.Process) (out *schema.Process, err error) {
	definition := p.extract(in.Definition.Input, in.Definition.Output, in.Definition.Error)
	result, err := p.store.UpdateProcess(in.Name, in.Version, definition)
	if err != nil {
		return
	}
	out = p.parse(result, in.Name)
	log.Printf("[Updated] %v", out)
	return
}

func (p *processServer) Activate(_ context.Context, in *schema.Options) (out *schema.Process, err error) {
	status := schema.Status_ACTIVATED.String()
	result, err := p.store.UpdateStatus(status, in.Name, in.Version)
	if err != nil {
		return
	}
	out = p.parse(result, in.Name)
	log.Printf("[Activated] %v", out)
	return
}

func (p *processServer) Deactivate(_ context.Context, in *schema.Options) (out *schema.Process, err error) {
	status := schema.Status_DEACTIVATED.String()
	result, err := p.store.UpdateStatus(status, in.Name, in.Version)
	if err != nil {
		return
	}
	out = p.parse(result, in.Name)
	log.Printf("[Deactivated] %v", out)
	return
}

func (p *processServer) extract(input, output, error *schema.Process_Data) (definition []ProcessData) {
	categories := []processCategory{processInput, processOutput, processError}
	processData := []*schema.Process_Data{input, output, error}
	for i, category := range categories {
		definition = append(definition, ProcessData{
			Kind:     "Event",
			Category: string(category),
			Name:     processData[i].Event.Name,
			Version:  processData[i].Event.Version,
		})
		if processData[i].Entity != nil {
			definition = append(definition, ProcessData{
				Kind:     "Entity",
				Category: string(category),
				Name:     processData[i].Entity.Name,
				Version:  processData[i].Entity.Version,
			})
		}
	}
	return
}

func (p *processServer) parse(process neo4j.Records, name string) (out *schema.Process) {
	status := schema.Status_UNACTIVATED
	if s, ok := process[name+".status"]; ok {
		status = schema.Status(schema.Status_value[s.(string)])
	}
	processData := make(map[string]*schema.Process_Data)
	for _, category := range []processCategory{processInput, processOutput, processError} {
		alias := fmt.Sprintf("%s_%s_", name, category)
		var event *schema.Event
		if name, ok := process[alias+"Event.name"]; ok {
			status := process[alias+"Event.status"].(string)
			event = &schema.Event{
				Name:    name.(string),
				Version: process[alias+"Event.version"].(string),
				Status:  schema.Status(schema.Status_value[status]),
			}
		}
		var entity *schema.Entity
		if name, ok := process[alias+"Entity.name"]; ok {
			status := process[alias+"Entity.status"].(string)
			entity = &schema.Entity{
				Name:    name.(string),
				Version: process[alias+"Entity.version"].(string),
				Status:  schema.Status(schema.Status_value[status]),
			}
		}
		processData[string(category)] = &schema.Process_Data{
			Event: event,
			Entity: entity,
		}
	}
	out = &schema.Process{
		Name:    process[name+".name"].(string),
		Version: process[name+".version"].(string),
		Status:  status,
		Definition: &schema.Process_Definition{
			Input:  processData["Input"],
			Output: processData["Output"],
			Error:  processData["Error"],
		},
	}
	return
}
