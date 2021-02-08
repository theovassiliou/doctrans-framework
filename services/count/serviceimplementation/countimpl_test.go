package serviceimplementation

import (
	"testing"

	"github.com/theovassiliou/doctrans-framework/dtaservice"
)

func TestExecuteWorkerLocally(t *testing.T) {
	type args struct {
		s        DtaService
		fileName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"simple",
			args{
				DtaService{
					IDocTransServer:              nil,
					GenDocTransServer:            dtaservice.GenDocTransServer{},
					UnimplementedDTAServerServer: dtaservice.UnimplementedDTAServerServer{},
				},
				"../../../test/testDoc.txt",
			},
			`{
  "Bytes": 55,
  "Lines": 3,
  "Words": 11
}`,
		},
		{
			"no valid file",
			args{
				DtaService{
					IDocTransServer:              nil,
					GenDocTransServer:            dtaservice.GenDocTransServer{},
					UnimplementedDTAServerServer: dtaservice.UnimplementedDTAServerServer{},
				},
				"../../../test/testDocInvalid.txt",
			},
			"",
		},
		{
			"empty file",
			args{
				DtaService{
					IDocTransServer:              nil,
					GenDocTransServer:            dtaservice.GenDocTransServer{},
					UnimplementedDTAServerServer: dtaservice.UnimplementedDTAServerServer{},
				},
				"",
			},
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExecuteWorkerLocally(tt.args.s, tt.args.fileName); got != tt.want {
				t.Errorf("ExecuteWorkerLocally() = %v, want %v", got, tt.want)
			}
		})
	}
}
