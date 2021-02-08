# Contents

This directory contains the DTA specification and a generic DTA server implementation

- [dtaservice.proto](proto/dtaservice.proto): The normative specification of the GRPC specification for the DTA services
- `*.pb.go`: Generated files. Created by `protoc` using `grpc-gateway`, `protoc-gen-govalidators` and `grpc-ecosystem`
- `*.go`: The reference implementation for a GRPC-HTTP accessible DTA implementation

NOTE: [../swagger/dtaservice.swagger.json](../swagger/dtaservice.swagger.json) contains a (generated) swagger specification for a REST-specification
of DTA.
