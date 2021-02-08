package dtaservice

import (
	"reflect"
	"testing"
	"time"

	"google.golang.org/grpc/metadata"
)

func TestGetXinstanceIdHeader(t *testing.T) {
	type args struct {
		s *GenDocTransServer
	}
	tests := []struct {
		name string
		args args
		want metadata.MD
	}{
		{
			"simple - Now()",
			args{
				&GenDocTransServer{
					XInstanceIDprefix:    "aa/bb/",
					XInstanceIDstartTime: time.Now(),
				},
			},
			metadata.Pairs("X-Instance-Id", "aa/bb/0s"),
		},
		{
			"simple - No startTime",
			args{
				&GenDocTransServer{
					XInstanceIDprefix: "aa/bb/",
				},
			},
			metadata.MD{},
		},
		{
			"No prefix - no time",
			args{
				&GenDocTransServer{},
			},
			metadata.MD{},
		},
		{
			"No prefix - with time",
			args{
				&GenDocTransServer{
					XInstanceIDstartTime: time.Now(),
				},
			},
			metadata.Pairs("X-Instance-Id", "0s"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetXinstanceIDHeader(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetXinstanceIdHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}
