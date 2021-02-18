package instanceid

import "testing"

func TestBuildVBC(t *testing.T) {
	type args struct {
		appName string
		version string
		branch  string
		commit  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "full",
			args: args{
				appName: "msA",
				version: "1.1",
				branch:  "development",
				commit:  "12345",
			},
			want: "msA/1.1/development-12345%",
		},
		{
			name: "full - no branch",
			args: args{
				appName: "msA",
				version: "1.1",
				branch:  "",
				commit:  "12345",
			},
			want: "msA/1.1/-12345%",
		},
		{
			name: "simple - no B, no C",
			args: args{
				appName: "msA",
				version: "1.1",
				branch:  "",
				commit:  "",
			},
			want: "msA/1.1%",
		},
		{
			name: "invalid - no V, no B, no C",
			args: args{
				appName: "msA",
				version: "",
				branch:  "",
				commit:  "",
			},
			want: "msA/%",
		},
		{
			name: "invalid empty",
			args: args{
				appName: "",
				version: "",
				branch:  "",
				commit:  "",
			},
			want: "/%",
		},
		{
			name: "valid - unicode name",
			args: args{
				appName: "LifeOf∏",
				version: "1.1",
				branch:  "",
				commit:  "",
			},
			want: "LifeOf∏/1.1%",
		},
		{
			name: "full - version only dotts",
			args: args{
				appName: "msA",
				version: "....",
				branch:  "development",
				commit:  "12345",
			},
			want: "msA/..../development-12345%",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildVBC(tt.args.appName, tt.args.version, tt.args.branch, tt.args.commit); got != tt.want {
				t.Errorf("BuildVBC() = %v, want %v", got, tt.want)
			}
		})
	}
}
