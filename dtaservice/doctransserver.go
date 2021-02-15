package dtaservice

import (
	"context"
	"net"
	"strconv"
	sync "sync"
	"time"

	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"

	"github.com/carlescere/scheduler"
	log "github.com/sirupsen/logrus"
	aux "github.com/theovassiliou/doctrans-framework/ipaux"
	"github.com/theovassiliou/go-eureka-client/eureka"
)

// DocTransServerOptions describes communication related options that a DTS offers to the user
type DocTransServerOptions struct {
	GRPC         bool   `opts:"group=Protocols" help:"Start service only with GRPC protocol support if set"`
	HTTP         bool   `opts:"group=Protocols" help:"Start service only with HTTP protocol support if set"`
	Port         int    `opts:"group=Protocols" help:"On which port (starting point) to listen for the supported protocol(s)."`
	XInstanceID  bool   `opts:"group=Protocols" help:"If set disable X-Instance-Id disclosure on request."` // my instance ID
	RegHostName  string `opts:"group=Service" help:"If provided will be used as hostname for registration, else automatically derived."`
	RegIPAddress string `opts:"group=Service" help:"If provided will be used as ip-address for registration, else automatically derived."`
	RegPort      string `opts:"group=Service" help:"If provided will be used as port for registration, else automatically derived."`
	RegistrarURL string `opts:"group=Registrar" help:"Registry URL (ex http://eureka:8761/eureka). If set to \"\", no registration to eureka"`
}

// DocTransServerGenericOptions describes generic options that a DTS server offers to the user
type DocTransServerGenericOptions struct {
	LogLevel log.Level `opts:"group=Generic" help:"Log level, one of panic, fatal, error, warn or warning, info, debug, trace"`
	CfgFile  string    `opts:"group=Generic" help:"The config file to use" json:"-"`
	Init     bool      `opts:"group=Generic" help:"Create a default config file as defined by cfg-file, if set. If not set ~/.dta/{AppName}/config.json will be created." json:"-"`
}

// IDocTransServer specifies the interfaces a DocTransServer must implement
type IDocTransServer interface {
	GetDocTransServer() GenDocTransServer
	DTAServerServer
}

// GenDocTransServer is the struct that contains a generic server.
type GenDocTransServer struct {
	AppName              string `opts:"-"`
	DtaType              string `opts:"-"`
	Proto                string `opts:"-"`
	XInstanceIDprefix    string `opts:"instance-id"` // my instance ID
	XInstanceIDstartTime time.Time

	registrar    *eureka.Client
	instanceInfo *eureka.InstanceInfo
	heartBeatJob *scheduler.Job
	// UnimplementedDTAServerServer
}

// StartGrpcServer starts the grpc server implementation for a given listener
func StartGrpcServer(lis net.Listener, dtaServer DTAServerServer) {
	s := grpc.NewServer()

	RegisterDTAServerServer(s, dtaServer)
	if err := s.Serve(lis); err != nil {
		log.WithFields(log.Fields{"Service": "Registrar", "Status": "Abort"}).Fatalf("failed to serve: %v", err)
	}

}

// LaunchServices starts the grpcGateway and/or the httpGateway, registers, if indicated the respectiv services at the provided eureka service, and
// enables registration/deregistration or toggling of thereof via signals (CTRL-D, CTRL-C)
func LaunchServices(grpcGateway, httpGateway IDocTransServer, appName, dtaType, homepageURL string, options DocTransServerOptions) {
	var gDTS, hDTS GenDocTransServer

	var _httpListener net.Listener
	var _httpPort int

	// create GRPC Listener
	// -- take initial port
	_initialPort := options.Port
	// -- start listener and save used grpc port
	_grpcListener, _grpcPort := CreateListener(_initialPort, 20)
	var _ipAddressUsed string
	if options.RegIPAddress != "" {
		_ipAddressUsed = options.RegIPAddress
	} else {
		_ipAddressUsed, _ = aux.ExternalIP()
	}
	var registerGRPC, registerHTTP bool

	var theGPort = _grpcPort
	if grpcGateway != nil {
		registerGRPC = true
		gDTS = grpcGateway.GetDocTransServer()

		if options.RegPort != "" {
			theGPort, _ = strconv.Atoi(options.RegPort)
		}
		gDTS.NewInstanceInfo(calcHostName("grpc", options.RegHostName), appName, _ipAddressUsed, theGPort,
			0, false, dtaType, "grpc",
			homepageURL,
			"",
			"")
	}

	if httpGateway != nil {
		registerHTTP = true
		// create HTTP Listener (optional)
		// -- take GRPC port + 1
		// -- start listener and save used http port
		_httpListener, _httpPort = CreateListener(_grpcPort+1, 20)
		hDTS = httpGateway.GetDocTransServer()
		var theHPort = _httpPort
		if options.RegPort != "" {
			theHPort = theGPort + (_httpPort - _grpcPort)
		}
		hDTS.NewInstanceInfo(calcHostName("http", options.RegHostName), appName, _ipAddressUsed, theHPort,
			0, false, dtaType, "http",
			homepageURL,
			calcStatusURL(options.RegHostName+":"+strconv.Itoa(_httpPort)),
			calcHealthURL(options.RegHostName+":"+strconv.Itoa(_httpPort)))
	}

	var wg sync.WaitGroup

	// Register at registrar
	// -- Register service with GRPC protocol
	log.Tracef("RegistrarURL: %s\n", options.RegistrarURL)
	if registerGRPC && options.RegistrarURL != "" {
		gDTS.RegisterAtRegistry(options.RegistrarURL)
	}
	if registerGRPC {
		go StartGrpcServer(_grpcListener, grpcGateway)
		CaptureSignals(grpcGateway, options.RegistrarURL, &wg)
		wg.Add(1)
	}

	// -- Register service with HTTP protocol (optional)
	if registerHTTP && options.RegistrarURL != "" {
		hDTS.RegisterAtRegistry(options.RegistrarURL)
	}

	if registerHTTP {
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		go MuxHTTPIntoGrpc(ctx, _httpListener, _grpcPort)
		CaptureSignals(httpGateway, options.RegistrarURL, &wg)
		wg.Add(1)
	}

	wg.Wait()
}

// ----- Default implementations of DTA functions

// ListServices standard implementation of listing the supported services
func (dtas *GenDocTransServer) ListServices(ctx context.Context, req *empty.Empty) (*ListServicesResponse, error) {
	log.WithFields(log.Fields{"Service": dtas.ApplicationName(), "Status": "ListServices"}).Tracef("Service requested")
	log.WithFields(log.Fields{"Service": dtas.ApplicationName(), "Status": "ListServices"}).Infof("In know only myself: %s", dtas.ApplicationName())
	services := (&ListServicesResponse{}).Services
	services = append(services, dtas.ApplicationName())
	return &ListServicesResponse{Services: services}, nil
}

// Options has no default implementation. Overwrite to provide own implementation
func (*GenDocTransServer) Options(context.Context, *empty.Empty) (*OptionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Options not implemented")
}

// TransformDocument has no default implementation. Overwrite to provide own implementation
func (*GenDocTransServer) TransformDocument(context.Context, *TransformDocumentRequest) (*TransformDocumentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TransformDocument not implemented")
}

// TransformDocumentPipe has no default implementation. Overwrite to provide own implementation
func (*GenDocTransServer) TransformDocumentPipe(context.Context, *TransformDocumentPipeRequest) (*TransformDocumentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TransformDocumentPipe not implemented")
}
