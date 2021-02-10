# InstanceID DTA Service

The InstanceID service translates an instanceID into a visual tree representation like

```text
.
└── [88s]  MsA/1.1
    ├── [5555s]  msC/1.4
    └── [23234s]  msD/2.2
```

## Short Name

instanceid

## Full App Name

DE.TU-BERLIN.INSTANCEID

## Description

```text
  Usage: instanceid [options]

  Protocols options:
  --grpc, -g             Start service only with GRPC protocol support if set
  --http, -h             Start service only with HTTP protocol support if set
  --port, -p             On which port (starting point) to listen for the supported protocol(s).
                         (default 50000)
  --x-instance-id        If set disable X-Instance-Id disclosure on request.

  Service options:
  --host-name            If provided will be used as hostname, else automatically derived. (default
                         Theofaniss-iMac.fritz.box)

  Registrar options:
  --registrar-url, -r    Registry URL (ex http://eureka:8761/eureka). If set to "", no registration
                         to eureka (default http://eureka:8761/eureka)

  Generic options:
  --log-level, -l        Log level, one of panic, fatal, error, warn or warning, info, debug, trace
                         (default warning)
  --cfg-file, -c         The config file to use (default /Users/the/.dta/DE.TU-BERLIN.INSTANCEID/config.json)
  --init, -i             Create a default config file as defined by cfg-file, if set. If not set
                         ~/.dta/{AppName}/config.json will be created.

  Local Execution options:
  --local-execution, -x  If set, execute the service locally once and read from this file

  Options:
  --version, -v          display version
  --help                 display help

  Version:
    instanceid unknown (git: main 352e3bf)

  Read more:
    github.com/theovassiliou/doctrans
```

## Example usage

```shell
bin/count -x test/testDoc.txt
{
  "Bytes": 55,
  "Lines": 3,
  "Words": 11
}
```

## Example transformation

The endpoint `/v1/document/transform`  transforms the instanceId `MsA/1.1/dev-git22%88s(msC/1.4%5555s+msD/2.2%23234s)` into

```text
.
└── [88s]  MsA/1.1
    ├── [5555s]  msC/1.4
    └── [23234s]  msD/2.2
```

## Implemented Options

- GRPC-Service
- REST-Service
- X-Instance-Identification (-x-instance-id)
- Eureka registration (-r)
- Local execution (-x)

## Implemented endpoints

- `/v1/document/transform`
- `/v1/service/list`

## Implemented without functionality

- `/v1/document/transform-pipe`
- `/v1/service/options`

## Additional information

This service is an example service to play around with.
