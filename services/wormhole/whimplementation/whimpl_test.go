package whimplementation

import "testing"

func Test_shortenFQName(t *testing.T) {
	type args struct {
		fqName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Normal 1",
			args: args{"A.B.C"},
			want: "A.B",
		},
		{
			name: "Normal 2 - Final Dot",
			args: args{"A.B.C."},
			want: "A.B",
		},
		{
			name: "Long",
			args: args{"A.B.C.C.D"},
			want: "A.B.C.C",
		},
		{
			name: "One",
			args: args{"A"},
			want: "",
		},
		{
			name: "None",
			args: args{""},
			want: "",
		},
		{
			name: "Two Dots",
			args: args{"A..B"},
			want: "A",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shortenFQName(tt.args.fqName); got != tt.want {
				t.Errorf("shortenFQName() = %v, want %v", got, tt.want)
			}
		})
	}
}
