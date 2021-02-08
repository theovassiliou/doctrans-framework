# 8. Example simple doctrans implementations

Date: 2021-01-26

## Status

Accepted

## Context

Within the framework a reference implementation is given. Examples demonstrate the usage of the framework. Both as a server and a client.
As we are also following the reverse-proxy approach for providing simultaneously GRPC and HTTP access to each service, example on how to implement
this would be also beneficial

## Decision

We will provide example implementations for

- GRPC services
- GRPC & HTTP services (reverse-proxy)
- GRPC clients
- HTTP clients

## Consequences

By providing this examples the documentation becomse much more conciese. This examples should not be considered as production ready. Each example might focus on a different aspect of the implementation framework.
DTA implementation that should be considered production ready should be maintenained independendly from this repo.

### Risks

- That we focus too much on the "elegance" of the examples than on the educational task. We agree that each example implementation is a standalone implementation and reflects the educational mission at the point of creation.
- Different implementation strategies could irritate users on how to use the framework.
