package whimplementation

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/sirupsen/logrus"
	"github.com/theovassiliou/doctrans-framework/dtaservice"
	pb "github.com/theovassiliou/doctrans-framework/dtaservice"
	"github.com/theovassiliou/doctrans-framework/sympan"
	"github.com/theovassiliou/go-eureka-client/eureka"
	instanceid "github.com/theovassiliou/instanceidentification"
)

// Wormhole holds the infrastructure for performing the service
type Wormhole struct {
	pb.UnimplementedDTAServerServer
	pb.GenDocTransServer
	pb.IDocTransServer
	Scope string
}

// TransformDocument looks up the requested services via the resolver and forwards the request to the resolved service.
func (dtas *Wormhole) TransformDocument(ctx context.Context, in *pb.TransformDocumentRequest) (*pb.TransformDocumentResponse, error) {
	resolver := dtas.GetResolver()
	fqServiceName := in.GetServiceName()
	var theSelectedInstance eureka.InstanceInfo
	theSelectedInstance, err := sympan.WormholeResolveApplication(resolver, dtas.Scope, fqServiceName, "grpc", true)

	if err != nil {
		var theError []string

		if err.Error() == "not found" {
			theError = append(theError, "could not resolve service")
		} else {
			theError = append(theError, err.Error())
		}
		theError = append(theError, "Could not find service "+in.GetServiceName())
		return &pb.TransformDocumentResponse{
			Document:    []byte{},
			TransOutput: []string{},
			Error:       theError,
		}, nil
	}
	return forwardRequest(ctx, dtas, theSelectedInstance, in)
}

func forwardRequest(ctx context.Context, dtas *Wormhole, theSelectedInstance eureka.InstanceInfo, in *pb.TransformDocumentRequest) (*pb.TransformDocumentResponse, error) {
	conn, err := grpc.Dial(theSelectedInstance.IpAddr+":"+theSelectedInstance.Port.Port, grpc.WithInsecure())
	if err != nil {
		log.WithFields(log.Fields{"Service": dtas.AppName, "Status": "TransformDocument"}).Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewDTAServerClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	var reqHeader metadata.MD

	r, err := c.TransformDocument(ctx, in, grpc.Header(&reqHeader))
	if err != nil {
		log.WithFields(log.Fields{"Service": dtas.GenDocTransServer.AppName, "Status": "TransformDocument"}).Fatalf("Failed to transform: %s", err.Error())
	}
	log.WithFields(log.Fields{"Service": dtas.GenDocTransServer.AppName, "Status": "TransformDocumentResult"}).Tracef("%s\n", string(r.GetDocument()))

	if dtas.XInstanceIDprefix != "" {
		g := dtas.GetDocTransServer()
		myMiid := dtaservice.CreateMiid(g)
		ciidString := reqHeader.Get("X-Instance-Id")
		var ciids []instanceid.Ciid

		if len(ciidString) == 0 { // No X-Instance-Id provided
			c := instanceid.NewCiid(theSelectedInstance.App + "/na/%-1s")
			ciids = append(ciids, c)
		} else {
			for _, c := range ciidString {
				ciids = append(ciids, instanceid.NewCiid(c))
			}
		}
		myCiid := instanceid.Ciid{
			Miid:  myMiid,
			Ciids: ciids,
		}
		header := metadata.Pairs("X-Instance-Id", myCiid.String())
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
