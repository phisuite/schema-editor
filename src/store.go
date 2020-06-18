package main

import (
	"errors"
	"fmt"
	"github.com/phiskills/neo4j-client.go"
)

type store struct {
	client        *neo4j.Client
	kind          string
	defaultStatus string
}

type query struct {
	neo4j.Query
}

type Field struct {
	Name, Type, Category string
}

type ProcessData struct {
	Kind, Category, Name, Version string
}

func (s *store) CreateData(name, version string, fields []Field) (neo4j.Records, error) {
	records, err := s.client.Write(func(j neo4j.Job) (result neo4j.Result, err error) {
		none, err := s.find(name, version, j)
		if err != nil {
			return
		}
		if none != nil {
			err = s.errorMessage(name, version, "already exists")
			return
		}
		data := &neo4j.Node{
			Id:     name,
			Labels: []string{s.kind},
			Props:  neo4j.Records{
				"name":    name,
				"version": version,
				"status":  s.defaultStatus,
			},
		}
		query := s.newRequest()
		query.Query = query.Create(data)
		query = query.mergeFields(fields, data)
		records, err := j.Execute(query.Query)
		if err != nil {
			return
		}
		result = records
		return
	})
	if err != nil {
		return nil, err
	}
	return records[0], err
}

func (s *store) UpdateData(name, version string, fields []Field) (neo4j.Records, error) {
	records, err := s.client.Write(func(j neo4j.Job) (result neo4j.Result, err error) {
		existing, err := s.find(name, version, j)
		if err != nil {
			return
		}
		if existing == nil {
			err = s.errorMessage(name, version, "not found")
			return
		}
		if existing[name+".status"] != s.defaultStatus {
			err = s.errorMessage(name, version, "already activated")
			return
		}
		data := &neo4j.Node{
			Id:     name,
			Labels: []string{s.kind},
			Props:  neo4j.Records{
				"name":    name,
				"version": version,
			},
		}
		path := &neo4j.Path{
			Origin:       &neo4j.Node{Id: name},
			Relationship: &neo4j.Relationship{
				Id:        "has",
				Type:      "Has",
			},
			Destination:  &neo4j.Node{
				Labels: []string{"Field"},
			},
		}
		query := s.newRequest()
		query.Query = query.Match(data).Optional().Match(path)
		query.Query = query.Delete("has")
		query = query.mergeFields(fields, data)
		records, err := j.Execute(query.Query)
		if err != nil {
			return
		}
		result = records
		return
	})
	if err != nil {
		return nil, err
	}
	return records[0], err
}

func (s *store) CreateProcess(name, version string, definition []ProcessData) (neo4j.Records, error) {
	records, err := s.client.Write(func(j neo4j.Job) (result neo4j.Result, err error) {
		none, err := s.find(name, version, j)
		if err != nil {
			return
		}
		if none != nil {
			err = s.errorMessage(name, version, "already exists")
			return
		}
		process := &neo4j.Node{
			Id:     name,
			Labels: []string{s.kind},
			Props:  neo4j.Records{
				"name":    name,
				"version": version,
				"status":  s.defaultStatus,
			},
		}
		query := s.newRequest()
		query.Query = query.Create(process)
		query = query.mergeDefinition(definition, process)
		records, err := j.Execute(query)
		if err != nil {
			return
		}
		result = records
		return
	})
	if err != nil {
		return nil, err
	}
	return records[0], err
}

func (s *store) UpdateProcess(name, version string, definition []ProcessData) (neo4j.Records, error) {
	records, err := s.client.Write(func(j neo4j.Job) (result neo4j.Result, err error) {
		existing, err := s.find(name, version, j)
		if err != nil {
			return
		}
		if existing == nil {
			err = s.errorMessage(name, version, "not found")
			return
		}
		if existing[name+".status"] != s.defaultStatus {
			err = s.errorMessage(name, version, "already activated")
			return
		}
		process := &neo4j.Node{
			Id:     name,
			Labels: []string{s.kind},
			Props:  neo4j.Records{
				"name":    name,
				"version": version,
			},
		}
		path := &neo4j.Path{
			Origin:       &neo4j.Node{Id: name},
			Relationship: &neo4j.Relationship{Id: "rel"},
		}
		query := s.newRequest()
		query.Query = query.Match(process).Match(path).Delete("rel")
		query = query.mergeDefinition(definition, process)
		records, err := j.Execute(query)
		if err != nil {
			return
		}
		result = records
		return
	})
	if err != nil {
		return nil, err
	}
	return records[0], err
}

func (s *store) UpdateStatus(status, name, version string) (neo4j.Records, error) {
	records, err := s.client.Write(func(j neo4j.Job) (result neo4j.Result, err error) {
		existing, err := s.find(name, version, j)
		if err != nil {
			return
		}
		if existing == nil {
			err = s.errorMessage(name, version, "not found")
			return
		}
		if existing[name+".status"] == status {
			err = s.errorMessage(name, version, "status already up-to-date")
			return
		}
		schema := &neo4j.Node{
			Id:     name,
			Labels: []string{s.kind},
			Props:  neo4j.Records{
				"name":    name,
				"version": version,
			},
		}
		setters := neo4j.Records{
			"status": status,
		}
		properties := schema.Properties("name", "version", "status")
		query := s.newRequest()
		query.Query = query.Match(schema).Set(schema, setters)
		query.Query = query.Return(properties...)
		records, err := j.Execute(query)
		if err != nil {
			return
		}
		result = records
		return
	})
	if err != nil {
		return nil, err
	}
	return records[0], err
}

func (s *store) errorMessage(name, version, message string) error {
	err := fmt.Sprintf("%s %s:%s %s", s.kind, name, version, message)
	return errors.New(err)
}

func (s *store) find(name, version string, j neo4j.Job) (neo4j.Records, error) {
	schema := &neo4j.Node{
		Id:     name,
		Labels: []string{s.kind},
		Props:  neo4j.Records{
			"name":    name,
			"version": version,
		},
	}
	properties := schema.Properties("name", "version", "status")
	query := s.client.NewRequest()
	query = query.Match(schema).Return(properties...)
	records, err := j.Execute(query)
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, nil
	}
	return records[0], nil
}

func (s *store) newRequest() query {
	return query{Query: s.client.NewRequest()}
}

func (q query) mergeFields(fields []Field, data *neo4j.Node) query {
	properties := data.Properties("name", "version", "status")
	for _, field := range fields {
		payload := &neo4j.Node{
			Id:     field.Name,
			Labels: []string{"Field"},
			Props: neo4j.Records{
				"name": field.Name,
				"type": field.Type,
			},
		}
		relationship := &neo4j.Relationship{
			Id:        data.Id + "_" + field.Name,
			Type:      "Has",
			Props:     neo4j.Records{"category": field.Category},
			Direction: neo4j.FromOriginToDestination,
		}
		q.Query = q.Merge(payload).Merge(&neo4j.Path{
			Origin:       &neo4j.Node{Id: data.Id},
			Relationship: relationship,
			Destination:  &neo4j.Node{Id: payload.Id},
		})
		properties = append(properties, payload.Properties("name", "type")...)
		properties = append(properties, relationship.Property("category"))
	}
	q.Query = q.Return(properties...)
	return q
}

func (q query) mergeDefinition(definition []ProcessData, process *neo4j.Node) query {
	properties := process.Properties("name", "version", "status")
	for _, data := range definition {
		alias := fmt.Sprintf("%s_%s_%s", process.Id, data.Category, data.Kind)
		node := &neo4j.Node{
			Id:     alias,
			Labels: []string{data.Kind},
			Props:  neo4j.Records{
				"name":    data.Name,
				"version": data.Version,
			},
		}
		relationship := &neo4j.Relationship{
			Id:        "rel_" + alias,
			Type:      data.Category,
			Direction: neo4j.FromOriginToDestination,
		}
		path := &neo4j.Path{
			Origin:       &neo4j.Node{Id: process.Id},
			Relationship: relationship,
			Destination:  &neo4j.Node{Id: node.Id},
		}
		q.Query = q.Merge(node).Create(path)
		properties = append(properties, node.Properties("name", "version", "status")...)
	}
	q.Query = q.Return(properties...)
	return q
}
