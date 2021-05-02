package serviceimplementation

import (
	"context"
	"io/ioutil"

	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/sirupsen/logrus"
	"github.com/theovassiliou/doctrans-framework/dtaservice"
	pb "github.com/theovassiliou/doctrans-framework/dtaservice"
	instanceid "github.com/theovassiliou/instanceidentification"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
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

// TransformDocument implements dtaservice.DTAServer
func (s *DtaService) TransformDocument(ctx context.Context, in *pb.TransformDocumentRequest) (*pb.TransformDocumentResponse, error) {

	// md, _ := metadata.FromIncomingContext(ctx)
	// log.Warnf("%#v", md)

	l, sOut, sErr := Work(s, in.GetDocument(), in.GetOptions())
	var errorS []string
	if sErr != nil {
		errorS = []string{sErr.Error()}
	}

	log.WithFields(log.Fields{"Service": "count", "Status": "TransformDocument"}).Tracef("Received document: %s and has lines %s", string(in.GetDocument()), l)
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

// Work returns an encoded JSON object containing the
// bytes 	count the number of bytes
// lines	count the numnber of lines
// words		count the number of words
// The Service returns  the number of lines, words, and bytes contained in the input document
func Work(s *DtaService, input []byte, options *structpb.Struct) (string, []string, error) {
	theID := string(input)

	parsed := instanceid.NewCiid(theID)

	// resB, err := json.MarshalIndent(parsed, "", " ")
	resB := instanceid.PrintCiid(parsed)
	log.WithFields(log.Fields{"Service": "Work"}).Debugf("The result %s", resB)

	return string(resB), []string{}, nil
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

	transDoc, _, err := Work(&s, dat, nil)
	if check(err) {
		return ""
	}

	return transDoc
}

// InstanceIdString executes the worker locally, and returns the transformation result as string
func InstanceIdString(s DtaService, iids string) string {
	if iids == "" {
		log.Errorln("No instanceId provided. Aborting.")
		return ""
	}

	transDoc, _, err := Work(&s, []byte(iids), nil)
	if check(err) {
		return ""
	}

	return transDoc
}

func MiidContained(s DtaService, iids, miid string) bool {
	if iids == "" {
		log.Errorln("No instanceId provided. Aborting.")
		return false
	}

	m := instanceid.NewCiid(iids)
	return m.Contains(miid)
}
