# 11. Version Control created files

Date: 2021-02-09

## Status

Accepted

## Context

When distributing a project, we have to consider whether we expect each repo user to build
all generated files. Not version controlling them means that each user builds from scratch, using the installed version
of tools. However, in the case of cloning the repo, the project does not compile, shows a lot of errors, and prohibits
the discussion of certain API artifacts, for example, the swagger files.

## Decision

We will version control generated source code files.

## Consequences

Building the repo is much easier, as `make build` immediately succeeds

### Risks

If no or incompatible versions of required tools are installed, this might be noticed late in the build process.
