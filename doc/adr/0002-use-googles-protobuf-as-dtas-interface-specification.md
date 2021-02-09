# 2. Use Googles protobuf as DTAs interface specification

Date: 2021-01-25

## Status

Accepted

Amended by [3. Use GRPC as main transport communication](0003-use-grpc-as-main-transport-communication.md)

## Context

We are defining a framework to create a series of various micro-services. For the API specification, we required an interface specification language.

We considered swagger/openAPI and protobuf as the primary API specification language.

## Decision

For transport efficiency reasons, we decided to use protbuf as the primary API interface specification language.

## Consequences

- We maintain the API specification in [dtaservice.proto](/dtaservice/dtaservice.proto)
- We require the use of [`protoc`](https://github.com/protocolbuffers/protobuf/releases)

### Risk

- In case we decide to support Swagger we have to maintain synchronicity
