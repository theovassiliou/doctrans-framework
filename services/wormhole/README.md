# Wormhole in Sympan

Wormholes in the Sympan architecture are entities that enable access to services in an inner galaxy from an outer galaxy. Other architectures are using the term gateway for similar functionality.

## Wormhole properties

A wormhole receives a service request and has to forward it into the galaxy it service.

Thus, a wormhole has to

- register itself to the outer-cadaster, so that it can be reached from the outer galaxy
- use the inner-cadaster , where he can look-up which stars (or galaxies) are present in the inner, the served galaxy
- register with his SCOPING-NAME and as service name use WH. Example `DE.TU-BERLIN.WH`

## Wormholes and IP Addresses

A wormhole has two facets. The one facets represents the wormhole in the outer galaxy. The wormhole has to be accessible via an IP Address that is "known" to in the outer galaxy. In addition it has to identify itself, by stating

- that it is a wormhole and not star. For this the reserved service name `WH` shall be used.
- Which scope it is serving.

While the sympan namespace is not related to the DNS wormholes and stars have to be reachable via IP. 

### Outer facet

### Inner facet

### Implementation assumptions

- To address a star or wormhole the IP addresses in the instance-information at the eureka shall be used, usage of hostnames is not recommended.
- A wormhole shall register itself
- A wormhole requires access to a resolver
- If a wormhole shall be accessible from the public internet it shall use it's public ipaddress for registration at the eureka service

## Scoping

The sympan architecture supports hierachical scopes.

The sympan can contain stars and galaxies, while galaxies can also contain stars and galaxies.

Example:

```plain
.
└── ECHO **
├── // DE.TU-BERLIN 
│   ├── WH **
│   ├── COUNT **
│   ├── // DE.TU-BERLIN.QDS
│   │   ├── INSTANCEID **
│   │   └── HTML2TEXT **
└── // BERLIN.VASSILIOU
```

`.` represents the root or the full sympan. It has one service