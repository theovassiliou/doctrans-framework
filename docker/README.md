# Docker and DTA services

DTA services can easily be deployed in Docker containers. We distinguish three different scenarios of how DTA services could be operated from within containers.

1. Service in a shared namespace
2. Service in a separate namespace
3. Service in an own galaxy

## Shared namespace

In this scenario, a single service is merely started in a container, is accessible via its port, and could be registered with its external addresses at an external eureka service.

### How to configure the Dockerfile.* file

- Create docker template files

`make docker-templates`

- Copy the desired template file to the docker directory. You should have something like

```Dockerfile
FROM scratch
EXPOSE 50000
ADD "./bin/docker/count" /
CMD ["/count", "-a", "count", "--reg-host-name", "localhost", "--reg-ip-address", "192.168.178.60","--reg-port", "60000","-g", "--registrar-url", "http://192.168.178.60:8761/eureka"]
```

As the service is living with its own IP-addresses and ports you can configure the information that is passed to the eureka server manually.

The above file assumes that you will export the exposed port ("50000") as port `60000`

Build your image with

`docker build -t count_grpc -f docker/Dockerfile.count.grpc .`

Run your images with
`docker run -p 60000:50000 -d -it count_grpc`

As a result, you will have a DTA service, that

- lives in a docker container
- has modified the service name to `counter`
- exposes the grpc port 50000 at 60000
- is registered at a eureka server at <http://192.168.178.60:8761/eureka>
- announces its availability at 192.168.178.60 and port 60000

All in all, while this works it requires quite a lot of manual work.

## Separate namespace

If you would like to run your service within its namespace you could modify via the `-a` parameter its name for example to `-a ORG.ACME.COUNT`.

This name will be then registered at the Eureka server. It is up to the client to process this information correctly.

An automatic namespace handling is introduced by using galaxies and wormholes.

## Own galaxy

An own galaxy consists of at least

- a wormhole (aka gateway)
- an own registry (aka eureka service)
- one or multiple stars (aka services)

To create a galaxy in a box we will use docker-compose to describe our galaxy.

Our galaxy will have the following properties:

- Namespace: DE.TU-BERLIN
- Galaxy Internal: GRPC communication only
- Registry: As registry we are using an EUREKA service
- Registry: Exposes Port for monitoring purpose
- Stars are reachable ONLY via a wormhole
- The wormwhole offers GRPC and HTTP services
- The wormhole registeres with an global registry

So let's compose the docker file:

```yaml
version: '3'
services:
  eureka:
     image: aista/eureka
     container_name: eureka-intern
     restart: always
     ports:
      - "9761:8761"
```

This instantiates an internal eureka server. For monitoring purposes we are exposing the internal port `8761` to `9761`. In an production environment we could remove this.

```yaml
  count_grpc:
    depends_on: 
      - eureka
    container_name: count_grpc
    restart: always
    build:
      context: ../
      dockerfile: docker/Dockerfile.count.grpc
    entrypoint:
    - /count
    - -g
    - -a
    - count
    - --registrar-url
    - http://eureka:8761/eureka
```
