# 6. Determine the public IP address

Date: 2021-01-26

## Status

Accepted

## Context

In order to be able to announce a publicly accessible IP address for service the DTA reference implementation requires some means to achieve this.
The only reliable approach that we found was using an external service.

## Decision

We are using an external service that can be accessed via HTTP. Initially, we will use  <https://api.ipify.org>

## Consequences

Communication to a service when starting a reference implementation. This might be not wanted.
