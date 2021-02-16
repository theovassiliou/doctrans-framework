package whimplementation

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/sirupsen/logrus"
	"github.com/theovassiliou/doctrans-framework/dtaservice"
	pb "github.com/theovassiliou/doctrans-framework/dtaservice"
)

// Wormhole holds the infrastructure for performing the service
type Wormhole struct {
	pb.UnimplementedDTAServerServer
	pb.GenDocTransServer
	pb.IDocTransServer
}

// TransformDocument looks up the requested services via the resolver and forwards the request to the resolved service.
func (dtas *Wormhole) TransformDocument(ctx context.Context, in *pb.TransformDocumentRequest) (*pb.TransformDocumentResponse, error) {

	// from message: which application is requested? (fully qualified service name)
	//  fqServiceName := fromMessage

	// from resolver: look for application
	// 		applicationExist := resolver.GetApplication(fqServiceName)

	// if available forward request to one instance
	// 		if applicationExist then
	//				theSelectedInstance := selectoOneOf(applicationExist)
	//				repsonse := theSelectedInstance.TransformDocument
	//				if response == successfull then
	//					return response, error
	//				else
	//					return _, error
	//
	// if not available (else) look for a wormhole that promises to implement it, by
	//		- find all wh's
	//		- seperate WHdomain from fqWHname and seperate domain from fqServiceNameDomain
	//		- LB:1 compare the fqServiceNameDomain with all WHdomains
	//		- if there is a full match
	//			- append wh again to matched WHdomain
	//			- get an instance
	//			- forward the request to the instance and return the result
	//		- if there is not a full match (else)
	// 			- remove the last element from fqServiceNameDomain
	//			- If there is something left: GOTO LB:1
	//			- Else return error

	// Let's find out whether we find the server that can serve this service.
	a, err := dtas.GetResolver().GetApplication(in.GetServiceName())
	if err != nil || len(a.Instances) == 0 {
		log.Errorf("Couldn't find server for app %s", in.GetServiceName())
		return &pb.TransformDocumentResponse{
			Document:    []byte{},
			TransOutput: []string{},
			Error:       []string{"Could not find service", "Service requested: " + in.GetServiceName()},
		}, nil
	}
	log.WithFields(log.Fields{"Service": dtas.AppName, "Status": "TransformDocument"}).Debugf("Connecting to: %s:%s", a.Instances[0].IpAddr, a.Instances[0].Port.Port)
	conn, err := grpc.Dial(a.Instances[0].IpAddr+":"+a.Instances[0].Port.Port, grpc.WithInsecure())
	if err != nil {
		log.WithFields(log.Fields{"Service": dtas.AppName, "Status": "TransformDocument"}).Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewDTAServerClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	r, err := c.TransformDocument(ctx, in)
	if err != nil {
		log.WithFields(log.Fields{"Service": dtas.GenDocTransServer.AppName, "Status": "TransformDocument"}).Fatalf("Failed to transform: %s", err.Error())
	}
	log.WithFields(log.Fields{"Service": dtas.GenDocTransServer.AppName, "Status": "TransformDocumentResult"}).Tracef("%s\n", string(r.GetDocument()))

	if dtas.XInstanceIDprefix != "" {
		g := dtas.GetDocTransServer()
		header := dtaservice.GetXinstanceIDHeader(g)
		grpc.SendHeader(ctx, header)
	}
	return r, err
}

// ListServices returns all the services visible for this gateway via the resolver
func (dtas *Wormhole) ListServices(ctx context.Context, req *empty.Empty) (*pb.ListServicesResponse, error) {
	// ListServices implements dtaservice.DTAServer
	log.Println(dtas.GetResolver())
	a, _ := dtas.GetResolver().GetApplications()

	log.WithFields(log.Fields{"Service": dtas.GenDocTransServer.AppName, "Status": "ListServices"}).Tracef("Service requested")
	log.WithFields(log.Fields{"Service": dtas.GenDocTransServer.AppName, "Status": "ListServices"}).Infof("Known Services registered with EUREKA: %v", a)
	services := (&pb.ListServicesResponse{}).Services
	for _, s := range a.Applications {
		services = append(services, s.Name)
	}

	if dtas.XInstanceIDprefix != "" {
		g := dtas.GetDocTransServer()
		header := dtaservice.GetXinstanceIDHeader(g)
		grpc.SendHeader(ctx, header)
	}
	return &pb.ListServicesResponse{Services: services}, nil
}

func (*Wormhole) TransformDocumentPipe(context.Context, *pb.TransformDocumentPipeRequest) (*pb.TransformDocumentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TransformPipe not implemented")
}
func (*Wormhole) Options(ctx context.Context, req *empty.Empty) (*pb.OptionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Options not implemented")
}

func (s *Wormhole) GetDocTransServer() *pb.GenDocTransServer {
	return &s.GenDocTransServer
}
