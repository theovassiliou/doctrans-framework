package serviceimplementation

import (
	"context"
	"fmt"
	"io/ioutil"
	"regexp"

	"github.com/Knetic/govaluate"
	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/sirupsen/logrus"
	"github.com/theovassiliou/doctrans-framework/dtaservice"
	pb "github.com/theovassiliou/doctrans-framework/dtaservice"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

type DtaService struct {
	pb.IDocTransServer
	pb.GenDocTransServer
	pb.UnimplementedDTAServerServer
}

var re *regexp.Regexp = regexp.MustCompile(`[\S]+`)

// Work returns an encoded JSON object containing the
// bytes 	count the number of bytes
// lines	count the numnber of lines
// words		count the number of words
// The Service returns  the number of lines, words, and bytes contained in the input document
func Work(s *DtaService, input []byte, options *structpb.Struct) (string, []string, error) {
	expString := string(input)
	expression, err := govaluate.NewEvaluableExpression(expString)
	var result interface{}
	if expression != nil {
		result, err = expression.Evaluate(nil)
	}
	return fmt.Sprintf("%v", result), []string{}, err

}

// TransformDocument implements dtaservice.DTAServer
func (s *DtaService) TransformDocument(ctx context.Context, in *pb.TransformDocumentRequest) (*pb.TransformDocumentResponse, error) {

	// md, _ := metadata.FromIncomingContext(ctx)
	// log.Warnf("%#v", md)

	result, sOut, sErr := Work(s, in.GetDocument(), in.GetOptions())
	var errorS []string
	if sErr != nil {
		errorS = []string{sErr.Error()}
	}

	log.WithFields(log.Fields{"Service": "add", "Status": "TransformDocument"}).Tracef("Received document: %s ", string(in.GetDocument()))
	// create and send header
	if s.XInstanceIDprefix != "" {
		g := s.GetDocTransServer()
		header := dtaservice.GetXinstanceIDHeader(g)
		grpc.SendHeader(ctx, header)
	}
	log.WithFields(log.Fields{"Service": "add", "Status": "TransformDocument"}).Tracef("Result: %s ", result)

	return &pb.TransformDocumentResponse{
		Document:    []byte(result),
		TransOutput: sOut,
		Error:       errorS,
	}, nil
}

func (s *DtaService) ListServices(ctx context.Context, req *empty.Empty) (*pb.ListServicesResponse, error) {
	log.WithFields(log.Fields{"Service": s.ApplicationName(), "Status": "ListServices"}).Tracef("Service requested")
	log.WithFields(log.Fields{"Service": s.ApplicationName(), "Status": "ListServices"}).Infof("In know only myself: %s", s.ApplicationName())
	services := (&pb.ListServicesResponse{}).Services
	services = append(services, s.ApplicationName())
	if s.XInstanceIDprefix != "" {
		g := s.GetDocTransServer()
		header := dtaservice.GetXinstanceIDHeader(g)
		grpc.SendHeader(ctx, header)
	}
	return &pb.ListServicesResponse{Services: services}, nil
}

func (*DtaService) TransformDocumentPipe(context.Context, *pb.TransformDocumentPipeRequest) (*pb.TransformDocumentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TransformPipe not implemented")
}
func (*DtaService) Options(ctx context.Context, req *empty.Empty) (*pb.OptionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Options not implemented")
}

func (s *DtaService) GetDocTransServer() *pb.GenDocTransServer {
	return &s.GenDocTransServer
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

	log.WithFields(log.Fields{"Service": "add", "Status": "TransformDocument"}).Tracef("Received document: %v ", dat)

	return transDoc
}
