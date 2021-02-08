# 4. Support RESTHTTP communication via reverse proxy possibility

Date: 2021-01-26

## Status

Accepted

Amends [3. Use GRPC as main transport communication](0003-use-grpc-as-main-transport-communication.md)

## Context

GRPC and the go-tools provide the possibility to enable REST based communication via a reverse proxy approach. [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) enables this.

 It reads protobuf service definitions and generates a reverse-proxy server which 'translates a RESTful HTTP API into gRPC. This server is generated according to the google.api.http annotations in your service definitions. We are using this to create bi-protocol servers, that support both GRPC and REST (optional) communication. And to create the swagger specification out of the proto specification.

## Decision

Enable RESTful communication with our GRPC based DTA implementations

## Consequences

We add the dependency to [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)
