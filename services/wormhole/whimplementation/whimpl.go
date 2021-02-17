package whimplementation

import (
	"context"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/sirupsen/logrus"
	"github.com/theovassiliou/doctrans-framework/dtaservice"
	pb "github.com/theovassiliou/doctrans-framework/dtaservice"
	"github.com/theovassiliou/doctrans-framework/instanceid"
	"github.com/theovassiliou/go-eureka-client/eureka"
)

// Wormhole holds the infrastructure for performing the service
type Wormhole struct {
	pb.UnimplementedDTAServerServer
	pb.GenDocTransServer
	pb.IDocTransServer
}

// TransformDocument looks up the requested services via the resolver and forwards the request to the resolved service.
func (dtas *Wormhole) TransformDocument(ctx context.Context, in *pb.TransformDocumentRequest) (*pb.TransformDocumentResponse, error) {
	resolver := dtas.GetResolver()

	// from message: which application is requested? (fully qualified service name)
	//  fqServiceName := fromMessage
	fqServiceName := in.GetServiceName()
	var theSelectedInstance eureka.InstanceInfo
	cont := true

	for cont {
		// from resolver: look for application
		// 		applicationExist := resolver.GetApplication(fqServiceName)
		app, err := resolver.GetApplication(fqServiceName)
		log.WithFields(log.Fields{"Service": dtas.AppName, "Status": "TransformDocument"}).Debugf("looking for %s", fqServiceName)
		log.WithFields(log.Fields{"Service": dtas.AppName, "Status": "TransformDocument"}).Debugf("apps %s has %v instances", app.Name, len(app.Instances))
		if err != nil || app == nil || len(app.Instances) <= 0 {
			log.WithFields(log.Fields{"Service": dtas.AppName, "Status": "TransformDocument"}).Fatalf("could not connect to resolver: %v", err)
			return nil, err
		}
		if app != nil && len(app.Instances) > 0 {
			theSelectedInstance = selectOneOf(app.Instances)
			return forwardRequest(dtas, theSelectedInstance, ctx, in)
		}

		fqServiceName = shortenFQName(fqServiceName)
		if len(fqServiceName) > 0 {
			fqServiceName = fqServiceName + ".WH"
		} else {
			cont = false
		}
	}

	return nil, nil
}

func shortenFQName(fqName string) string {
	fqName = strings.TrimSuffix(fqName, ".")
	elements := strings.Split(fqName, ".")
	fqName = strings.Join(elements[:len(elements)-1], ".")
	fqName = strings.TrimSuffix(fqName, ".")
	return fqName
}

func forwardRequest(dtas *Wormhole, theSelectedInstance eureka.InstanceInfo, ctx context.Context, in *pb.TransformDocumentRequest) (*pb.TransformDocumentResponse, error) {
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
		log.Println(ciidString)
		var ciids []instanceid.Ciid

		for _, c := range ciidString {
			ciids = append(ciids, instanceid.NewCiid(c))
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

func selectOneOf(instances []eureka.InstanceInfo) eureka.InstanceInfo {

	// TODO: Find a better selection mechanism
	if len(instances) > 0 {
		for _, i := range instances {
			if strings.HasPrefix(i.HostName, "grpc@") {
				return i
			}
		}
	}
	return eureka.InstanceInfo{}
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
