# Sympan

The sympan framework (from greek: sympan == universe) supports the sympan name service resolution as described in !!!FIXME

## The Problem solved with Sympan

Consider a set of micro-services, that

    a) implement the same API
    b) are developed by different organizations
    c) have potentially the same (local) name
    d) should be accessible by various clients

As all services implement the same API (a) and should be accessible by various clients (d) a way to access each individual micro-service.

Micro-services with the same local name, have to be made globally distinguishable.

## Sympan Name Service

The Sympan Name Service (SNS) is a set of rules, functions, and data types to build a hierarchy of names to be able to uniquely identify and access micro-services of a given type.

### Elements of Sympan

To scope the naming of services we define

    - sympan or universe
    - galaxy

For managing the elements of one scoping level we are using

    - cadaster: maintains a list and addresses of all elements within its scope
    - cadaster client (CC): is a client that interacts with a cadaster to resolve names to addresses

Conceptionally, within an SNS system we need the following active elements

    - wormhole: a service that acts as a gateway into a galaxy 
    - star: the service itself

### Overview

The following picture shows a general SNS universe for a given service type

```
universe (for service type) 
  └── galaxyA
    └── galaxyA.1
        └── galaxy A.1.X
            └── star A.1.X.S1
            └── star A.1.X.S2
  └── galaxyB
    └── starB.S1
    └── starB.S2
  └── galaxyC
  └── star S1 
```

The given `universe` consists of 5 galaxies. Three galaxies (`galaxyA`, `galaxyB`, `galaxyC`) are top-level galaxies, i.e. they are embedded directly in `universe`.

However, within `galaxyA` there are two additional galaxies embedded, `galaxyA.1` (directly) and `galaxyA.1.X` (embedded in `galaxyA.1`).

While `universe`, `galaxyA.1.X` and `galaxyB`contain star `galaxyC`, `galaxyA.1` and `galaxyA` have no stars.


## Rational

While SNS resembles some of the elements of the Domain Name Service (DNS) we have decided to specify an own naming service, to avoid ambiguities.
