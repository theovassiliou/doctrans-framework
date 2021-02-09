# 3. Use GRPC as main transport communication

Date: 2021-01-25

## Status

Accepted

Amends [2. Use Googles protobuf as DTAs interface specification](0002-use-googles-protobuf-as-dtas-interface-specification.md)

Amended by [4. Support RESTHTTP communication via reverse proxy possibility](0004-support-resthttp-communication-via-reverse-proxy-possibility.md)

## Context

By using protoc as the main API language, we have to decide also for a transport protocol

## Decision

We are using GRPC as the main transport protocol

## Consequences

GRPC uses HTTP/2 as a carrier and sits on top.
