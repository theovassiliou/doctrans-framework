package echoimplementation

import (
	"context"
	"io/ioutil"
	"regexp"

	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/sirupsen/logrus"
	"github.com/theovassiliou/doctrans-framework/dtaservice"
	pb "github.com/theovassiliou/doctrans-framework/dtaservice"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DtaService represents the service as offered by the DTA server
type DtaService struct {
	pb.IDocTransServer
	pb.GenDocTransServer
	pb.UnimplementedDTAServerServer
}

// CountResults describes the results of the transformation
type CountResults struct {
	Bytes int
	Lines int
	Words int
}

var re *regexp.Regexp = regexp.MustCompile(`[\S]+`)

// TransformDocument implements dtaservice.DTAServer
func (s *DtaService) TransformDocument(ctx context.Context, in *pb.TransformDocumentRequest) (*pb.TransformDocumentResponse, error) {

	l, sOut, sErr := Work(in.GetDocument())
	var errorS []string
	if sErr != nil {
		errorS = []string{sErr.Error()}
	} else {
		errorS = []string{}
	}
	log.WithFields(log.Fields{"Service": s.ApplicationName(), "Status": "TransformDocument"}).Tracef("Received document: %s and echoing", string(in.GetDocument()))

	// create and send header
	if s.XInstanceIDprefix != "" {
		g := s.GetDocTransServer()
		header := dtaservice.GetXinstanceIDHeader(g)
		grpc.SendHeader(ctx, header)
	}

	return &pb.TransformDocumentResponse{
		Document:    []byte(l),
		TransOutput: sOut,
		Error:       errorS,
	}, nil
}

// ListServices lists the services, that this DTA instance is offering.
func (s *DtaService) ListServices(ctx context.Context, req *empty.Empty) (*pb.ListServicesResponse, error) {
	log.WithFields(log.Fields{"Service": s.ApplicationName(), "Status": "ListServices"}).Tracef("Service requested")
	log.WithFields(log.Fields{"Service": s.ApplicationName(), "Status": "ListServices"}).Infof("In know only myself: %s", s.ApplicationName())
	services := (&pb.ListServicesResponse{}).Services
	services = append(services, s.ApplicationName())
	return &pb.ListServicesResponse{Services: services}, nil
}

// TransformDocumentPipe is currently not implemented
func (*DtaService) TransformDocumentPipe(context.Context, *pb.TransformDocumentPipeRequest) (*pb.TransformDocumentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TransformPipe not implemented")
}

// Options is currently not implemented
func (*DtaService) Options(ctx context.Context, req *empty.Empty) (*pb.OptionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Options not implemented")
}

// GetDocTransServer returns the server instance of this service
func (s *DtaService) GetDocTransServer() *pb.GenDocTransServer {
	return &s.GenDocTransServer
}

// Work just retuns the document (ECHO)
func Work(input []byte) (string, []string, error) {
	return string(input), []string{}, nil
}

func check(e error) bool {
	if e != nil {
		log.Errorln(e)
		return true
	}
	return false
}

// ExecuteWorkerLocally executes the worker locally, and returns the transformation result as string
func ExecuteWorkerLocally(s DtaService, fileName string) string {
	if fileName == "" {
		log.Errorln("No fileName on local executing provided. Aborting.")
		return ""
	}

	dat, err := ioutil.ReadFile(fileName)
	if check(err) {
		return ""
	}

	transDoc, _, err := Work(dat)
	if check(err) {
		return ""
	}

	return transDoc
}
