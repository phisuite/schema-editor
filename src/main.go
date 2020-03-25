package main

import (
	"github.com/phiskills/grpc-api.go"
	"github.com/phisuite/schema.go"
)

func main() {
	api := grpc.New()
	schema.RegisterEventWriteAPIServer(api.Server, &eventServer{})
	schema.RegisterEntityWriteAPIServer(api.Server, &entityServer{})
	schema.RegisterProcessWriteAPIServer(api.Server, &processServer{})
	api.Start()
}
