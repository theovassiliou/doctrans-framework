package dtaservice

import (
	"strconv"
	"time"

	"github.com/theovassiliou/doctrans-framework/instanceid"
	"google.golang.org/grpc/metadata"
)

// This file contains functionality to implement identity disclosure via
// X-Instance-Header fields

// GetXinstanceIDHeader returns, if identifity disclouse is indicated a, grpc metada to
// respond with an appropriate header suitable to
func GetXinstanceIDHeader(s *GenDocTransServer) metadata.MD {
	var x time.Time
	if s.XInstanceIDstartTime != x {
		return metadata.Pairs("X-Instance-Id", CreateMiidString(s))
	}

	return metadata.MD{}
}

func CreateMiidString(s *GenDocTransServer) string {
	var x time.Time
	if s.XInstanceIDstartTime != x {
		return s.XInstanceIDprefix + strconv.Itoa(int(time.Since(s.XInstanceIDstartTime).Seconds())) + "s"
	}
	return ""
}

func CreateMiid(s *GenDocTransServer) instanceid.Miid {
	return instanceid.NewMiid(CreateMiidString(s))
}
