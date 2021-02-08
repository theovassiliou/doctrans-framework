# 10. Support of Instance-Identification approach

Date: 2021-01-28

## Status

Accepted

## Context

In a paper submitted to "12th IEEE International Conference on Software Testing, Verification and Validation (ICST) 2020" and "2nd ACM/IEEE International Conference on Automation of Software Test AST 2021 " we proposed the introduction of an instance identification via the X-Identification-Id header. Our reference implemenation framework should support the usage of X-Instance-Id header fields and demonstrate the applicability.

## Decision

Implementing the support of instance-identification via HTTP and GRPC Headers

- X-Instance-Id for HTTP and
- x-instance-id for GRPC

## Consequences

- With the functionality provided in the framework we can easily include support for instance identification in our example doctrans implementations
- We have the possibility to evaluate the performance aspects of the ID generation and bandwith consumption, potentially
