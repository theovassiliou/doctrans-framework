# 5. Inventory of used libraries and tools

Date: 2021-01-26

## Status

Accepted

Amended by [7. Inventory of external services used](0007-inventory-of-external-services-used.md)

## Context

DTA is meant to be a research and exploration framework for micro-services. Therefore we need to have low-key access to the specification and the implementations. Documentation of processes and used tools is a prerequisite for this. Good examples might be another one.

## Decision

We maint a list of tools required to build and to play around with the DTA framework.

## Consequences

- Create [TOOLS.md](../TOOLS.md) and capture the used tools
- Won't capture there any required libraries. We expect that the go development environment, in particular, `go mod` should be sufficient to capture these dependencies.
