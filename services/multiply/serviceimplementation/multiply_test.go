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
				"../../../test/add1.txt",
			},
			`8`,
		},
		{
			"simple",
			args{
				DtaService{
					IDocTransServer:              nil,
					GenDocTransServer:            dtaservice.GenDocTransServer{},
					UnimplementedDTAServerServer: dtaservice.UnimplementedDTAServerServer{},
				},
				"../../../test/add2.txt",
			},
			`973`,
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
		{
			"empty file",
			args{
				DtaService{
					IDocTransServer:              nil,
					GenDocTransServer:            dtaservice.GenDocTransServer{},
					UnimplementedDTAServerServer: dtaservice.UnimplementedDTAServerServer{},
				},
				"../../../test/emptyFile.txt",
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
