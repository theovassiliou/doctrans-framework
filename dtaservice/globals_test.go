package dtaservice

import "testing"

func TestFormatFullVersion(t *testing.T) {
	type args struct {
		cmdName string
		version string
		branch  string
		commit  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"simple",
			args{
				"a", "b", "c", "d",
			},
			"a b (git: c d)",
		},
		{
			"empty b",
			args{
				"a", "", "c", "d",
			},
			"a unknown (git: c d)",
		},
		{
			"empty c",
			args{
				"a", "b", "", "d",
			},
			"a b (git: unknown d)",
		},
		{
			"empty d",
			args{
				"a", "b", "c", "",
			},
			"a b (git: c unknown)",
		},
		{
			"all empty",
			args{
				"", "", "", "",
			},
			" unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatFullVersion(tt.args.cmdName, tt.args.version, tt.args.branch, tt.args.commit); got != tt.want {
				t.Errorf("FormatFullVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
