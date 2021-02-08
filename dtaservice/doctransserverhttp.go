package dtaservice

import (
	"context"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	grpc "google.golang.org/grpc"

	log "github.com/sirupsen/logrus"
)

// MuxHTTPIntoGrpc starts the HTTP server in a given context, with a given listener. grpcPort must contain the port number of the GRPC server
// in addition the http operations as defined by the GRPC reverse proxy to additional endpoint (/status and /health) are being registered
func MuxHTTPIntoGrpc(ctx context.Context, httpListener net.Listener, grpcPort int) {

	incomingHeaders := func(header string) (string, bool) {
		if header == "X-Instance-Id" {
			return "x-instance-id", true
		}
		return runtime.DefaultHeaderMatcher(header)
	}

	outgoingHeaders := func(header string) (string, bool) {
		if header == "x-instance-id" {
			return "X-Instance-Id", true
		}
		return runtime.DefaultHeaderMatcher(header)
	}

	rmuxOptions1 := runtime.WithIncomingHeaderMatcher(incomingHeaders)
	rmuxOptions2 := runtime.WithOutgoingHeaderMatcher(outgoingHeaders)

	gwmux := runtime.NewServeMux(rmuxOptions1, rmuxOptions2)

	opts := []grpc.DialOption{grpc.WithInsecure()}
	log.WithFields(log.Fields{"Service": "HTTP"}).Debugf("GRPC Endpoint localhost:%d\n", grpcPort)

	err := RegisterDTAServerHandlerFromEndpoint(ctx, gwmux, "localhost:"+strconv.Itoa(grpcPort), opts)
	if err != nil {
		log.WithFields(log.Fields{"Service": "HTTP", "Status": "Abort"}).Fatalf("failed to register: %v", err)
	}

	// FIXME: Continue here and pull the handler out. Remember to change this in all services
	mux := http.NewServeMux()
	mux.HandleFunc("/status", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("X-Instance-Id", "some id")
		_, _ = io.Copy(w, strings.NewReader("The service is alive"))
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("X-Instance-Id", "some id")
		_, _ = io.Copy(w, strings.NewReader("The service is healthy as it is responding."))
	})

	mux.Handle("/", gwmux)

	// (4) Start HTTP Server
	// Start HTTP server (and proxy calls to gRPC server endpoint)
	log.WithFields(log.Fields{"Service": "HTTP", "Status": "Running"}).Debugf("Starting HTTP server on: %v", httpListener.Addr().String())

	if err := http.Serve(httpListener, mux); err != nil {
		log.WithFields(log.Fields{"Service": "HTTP", "Status": "Abort"}).Fatalf("failed to serve: %v", err)
	}
}

func calcStatusURL(instanceID string) string {
	return "http://" + instanceID + "/status"
}

func calcHealthURL(instanceID string) string {
	return "http://" + instanceID + "/health"
}
