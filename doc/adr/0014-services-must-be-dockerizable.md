# 14. services must be dockerizable

Date: 2021-02-15

## Status

Accepted

## Context

With quite a few services created we need to collect multiple services in a container so that the management is eased.
If multiple services live in one container we assume that the container will spawn its galaxy.

## Decision

1. Create makefiles to create flexible images for each or multiple services
2. If multiple-services live in one docker it has its namespace
3. A container with an own namespace will operate its own eureka service
4. A container with an own namespace will operate its gateway

## Consequences

By having separate images for each service, the deployment of each service will get easier.
We will be able to create more complex research environments, as we can support galaxies with separate namespaces
