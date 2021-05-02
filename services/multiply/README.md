# Multiply DTA Service

The multiply service multiplies two numbers, separated by "*".

Example input

```text
3 * 5
```

Example result

```text
15
````


## Short Name

multiply

## Full App Name

MATH.MULTIPLY

## Description

```text
  Usage: count [options]

  Protocols options:
  --grpc, -g             Start service only with GRPC protocol support if set
  --http, -h             Start service only with HTTP protocol support if set
  --port, -p             On which port (starting point) to listen for the supported protocol(s).
  --x-instance-id        Support X-Instance-Id disclosure on request.

  Service options:
  --host-name            If provided will be used as hostname, else automatically derived. (default {system dependent})

  Registrar options:
  --registrar-url, -r    Registry URL (ex http://eureka:8761/eureka, default http://eureka:8761/eureka)

  Generic options:
  --log-level, -l        Log level, one of panic, fatal, error, warn or warning, info, debug, trace
                         (default warning)
  --cfg-file, -c         The config file to use (default /Users/the/.dta/DE.TU-BERLIN.QDS.COUNT/config.json)
  --init, -i             Create a default config file as defined by cfg-file, if set. If not set
                         ~/.dta/{AppName}/config.json will be created.

  Local Execution options:
  --local-execution, -x  If set, execute the service locally once and read from this file

  Options:
  --version, -v          display version
  --help                 display help
```

## Example usage

```shell
bin/multiply -x "3*4"
12

bin/multiply -x "3*-4"
-12

```

## Example transformation

The endpoint `/v1/document/transform`  transforms the document 

```text
3*-4
```

into

```text
-12
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
- `/v1/document/transform-pipe`
- `/v1/service/options`

## Additional information

This service is an example service to play around
