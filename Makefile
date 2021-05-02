# list of servers to build
servers := "./services/count" "./services/instanceid" "./services/multiservices" "./services/wormhole"
clients := "./client/eurekabrowser"
gofiles := $(subst ./services/,,$(servers))
dockerexes := $(subst ./services/,./docker/,$(servers))
dockerfiles := $(subst ./services/,./docker/Dockerfile.,$(servers))
services := $(subst ./services/,,$(servers))


BINARY = doctrans-framework


VERSION?=?

ifeq ($(OS), Windows_NT)
        VERSION := $(shell git describe --exact-match --tags 2>nil)
else
        VERSION := $(shell git describe --exact-match --tags 2>/dev/null)
endif

COMMIT=$(shell git rev-parse --short HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

# Symlink into GOPATH
GITHUB_USERNAME=theovassiliou
BUILD_DIR=${GOPATH}/src/github.com/${GITHUB_USERNAME}/${BINARY}
BIN_DIR=${BUILD_DIR}/bin
CURRENT_DIR=$(shell pwd)
BUILD_DIR_LINK=$(shell readlink ${BUILD_DIR})


# list of generated files
genfiles := ./dtaservice/*.pb.* "./swagger/dtaservice.swagger.json"

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS = -ldflags "-X main.version=${VERSION} -X main.commit=${COMMIT} -X main.branch=${BRANCH}"

build: dtaservice/dtaservice.pb.go
	@mkdir -p $(BIN_DIR)
	go build ${LDFLAGS} -o ${BIN_DIR} ./...

./swagger/dtaservice.swagger.json dtaservice/dtaservice.pb.go: dtaservice/proto/dtaservice.proto
	protoc -Idtaservice/proto \
		-I/usr/local/include -I. \
		-I$(GOPATH)/src \
		-I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--go_out=plugins=grpc:dtaservice \
  		--grpc-gateway_out=logtostderr=true:dtaservice \
  		--swagger_out=logtostderr=true:swagger \
		--govalidators_out=dtaservice \
	dtaservice.proto

clean: 
ifneq ($(genfiles),)
	rm -f $(genfiles)
endif
ifneq ($(dockerexes),)
	rm -f $(dockerexes)
endif
	rm -rf $(BIN_DIR)

# Build the project
all: clean build test-all vet

# ADR: 9. Tests for all framework functionality
.PHONY: test
test:
	go test -short ./...

# Use this to get an idea on how to use the API
.PHONY: api-examples
api-examples:
	go test -v -tags=api_examples ./...


.PHONY: fmt
fmt:
	@gofmt -s -w $(GOFILES)

.PHONY: fmtcheck
fmtcheck:
	@if [ ! -z "$(GOFMT)" ]; then \
		echo "[ERROR] gofmt has found errors in the following files:"  ; \
		echo "$(GOFMT)" ; \
		echo "" ;\
		echo "Run make fmt to fix them." ; \
		exit 1 ;\
	fi

.PHONY: vet
vet:
	@echo 'go vet $$(go list ./...)'
	@go vet $$(go list ./...) ; if [ $$? -ne 0 ]; then \
		echo ""; \
		echo "go vet has found suspicious constructs. Please remediate any reported errors"; \
		echo "to fix them before submitting code for review."; \
		exit 1; \
	fi

.PHONY: check
check: fmtcheck vet

.PHONY: test-all
test-all: fmtcheck vet
	go test ./...

# DOCKERFILE AND DOCKERSERVICES MANAGEMENT

docker: dockerexes

.PHONY: dockerexes

dockerexes:
	@mkdir -p $(BIN_DIR)/docker
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ${BIN_DIR}/docker ./services/...

# Creates template Dockerfile for building images
docker-templates: 
	for n in $(services); do \
		rm -f docker/templates/Dockerfile.$$n.grpc; \
		rm -f docker/templates/Dockerfile.$$n.html; \
		rm -f docker/templates/Dockerfile.$$n.grpc+html; \
		echo "FROM scratch" >> docker/templates/Dockerfile.$$n.grpc ; \
		echo "EXPOSE 50000" >> docker/templates/Dockerfile.$$n.grpc ; \
		echo "ADD \"./bin/docker/$$n"\" / >> docker/templates/Dockerfile.$$n.grpc ; \
		echo "CMD [\"/$$n\", \"-a\", \"$$n\", \"--reg-host-name\", \"localhost\", \"--reg-ip-address\", \"192.168.178.60\",\"--reg-port\", \"60000\",\"-g\", \"--registrar-url\", \"http://192.168.178.60:8761/eureka\"]" >> docker/templates/Dockerfile.$$n.grpc ; \
		echo "FROM scratch" >> docker/templates/Dockerfile.$$n.html ; \
		echo "EXPOSE 50000" >> docker/templates/Dockerfile.$$n.html ; \
		echo "ADD \"./bin/docker/$$n"\" / >> docker/templates/Dockerfile.$$n.html ; \
		echo "CMD [\"/$$n\", \"-a\", \"$$n\", \"--reg-host-name\", \"localhost\", \"--reg-ip-address\", \"192.168.178.60\",\"--reg-port\", \"60000\",\"-h\", \"--registrar-url\", \"http://192.168.178.60:8761/eureka\"]" >> docker/templates/Dockerfile.$$n.html ; \
		echo "FROM scratch" >> docker/templates/Dockerfile.$$n.grpc+html ; \
		echo "EXPOSE 50000" >> docker/templates/Dockerfile.$$n.grpc+html ; \
		echo "ADD \"./bin/docker/$$n"\" / >> docker/templates/Dockerfile.$$n.grpc+html ; \
		echo "CMD [\"/$$n\", \"-a\", \"$$n\", \"--reg-host-name\", \"localhost\", \"--reg-ip-address\", \"192.168.178.60\",\"--reg-port\", \"60000\", \"--registrar-url\", \"http://192.168.178.60:8761/eureka\"]" >> docker/templates/Dockerfile.$$n.grpc+html ; \
	done


# with this command you can print the value of a Makefile variable
# Example: make print-genfiles
print-%  : ; @echo $* = $($*)