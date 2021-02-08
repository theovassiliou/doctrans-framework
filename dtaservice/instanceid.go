package dtaservice

import (
	"strconv"
	"time"

	"google.golang.org/grpc/metadata"
)

// This file contains functionality to implement identity disclosure via
// X-Instance-Header fields

// GetXinstanceIDHeader returns, if identifity disclouse is indicated a, grpc metada to
// respond with an appropriate header suitable to
func GetXinstanceIDHeader(s *GenDocTransServer) metadata.MD {
	var x time.Time
	if s.XInstanceIDstartTime != x {
		return metadata.Pairs("X-Instance-Id", s.XInstanceIDprefix+strconv.Itoa(int(time.Since(s.XInstanceIDstartTime).Seconds()))+"s")
	}

	return metadata.MD{}
}
