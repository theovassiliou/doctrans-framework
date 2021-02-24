# Multiservice DTA Service

The multiservice serves as an example for implementing multiple DTA services in one executable.
Multiservice implements an `Echo` and an `HTML2Text` service. For testing purposes, each service can be called locally. 

- Echo service: The echo service replies to the content that has been sent to it
- HTML2Text service: The HTML2Text service translates an HTML document passed to it to its textual representation, link, and tabular information is tried to be preserved.

## Short Name

multiservice

## Full App Name

MULTISERVICE

## Description

```text
  Usage: multiservices [options]

  Protocols options:
  --grpc, -g             Start service only with GRPC protocol support if set
  --http, -h             Start service only with HTTP protocol support if set
  --port, -p             On which port (starting point) to listen for the supported protocol(s).
                         (default 50000)
  --x-instance-id        If set disable X-Instance-Id disclosure on request.

  Service options:
  --reg-host-name, -r    If provided will be used as hostname for registration, else automatically
                         derived. (default Theofaniss-iMac.fritz.box)
  --reg-ip-address       If provided will be used as ip-address for registration, else automatically
                         derived.
  --reg-port             If provided will be used as port for registration, else automatically
                         derived.

  Registrar options:
  --registrar-url        Registry URL (ex http://eureka:8761/eureka). If set to "", no registration
                         to eureka (default http://eureka:8761/eureka)

  Generic options:
  --log-level, -l        Log level, one of panic, fatal, error, warn or warning, info, debug, trace
                         (default warning)
  --cfg-file, -c         The config file to use (default /Users/the/.dta/multiservice/config.json)
  --init, -i             Create a default config file as defined by cfg-file, if set. If not set
                         ~/.dta/{AppName}/config.json will be created.

  Local Execution options:
  --local-execution, -x  If set, execute the service locally once and read from this file
  --htm-l2-text, -1      If set, use HTML2TEXT service
  --echo, -2             If set, use ECHO service

  Options:
  --version, -v          display version
  --help                 display help

  Version:
    multiservice unknown (git: dockerize f4d99cd)

  Read more:
    github.com/theovassiliou/doctrans
```

## Echo Example usage

```shell
bin/multiservice -2 -x test/testDoc.txt
This is a test file
It has multiple lines
Basically 3.
```

## HTML2Text Example usage

```shell
bin/multiservice -1 -x test/html1.txt
Mega Service ( http://jaytaylor.com/ )

******************************************
Welcome to your new account on my service!
******************************************

Here is some more information:

* Link 1: Example.com ( https://example.com )
* Link 2: Example2.com ( https://example2.com )
* Something else

+-------------+-------------+
|  HEADER 1   |  HEADER 2   |
+-------------+-------------+
| Row 1 Col 1 | Row 1 Col 2 |
| Row 2 Col 1 | Row 2 Col 2 |
+-------------+-------------+
|  FOOTER 1   |  FOOTER 2   |
+-------------+-------------+
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
