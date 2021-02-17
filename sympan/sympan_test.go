package sympan

import (
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/theovassiliou/go-eureka-client/eureka"
)

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
			if got := ShortenFQName(tt.args.fqName); got != tt.want {
				t.Errorf("shortenFQName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_selectOneOf(t *testing.T) {
	type args struct {
		instances []eureka.InstanceInfo
		proto     string
	}
	tests := []struct {
		name string
		args args
		want eureka.InstanceInfo
	}{
		{
			name: "short",
			args: args{
				instances: []eureka.InstanceInfo{
					{
						HostName: "grpc@A",
						App:      "A",
					},
				},
				proto: "grpc",
			},
			want: eureka.InstanceInfo{
				HostName: "grpc@A",
				App:      "A",
			},
		},
		{
			name: "short - no grpc",
			args: args{
				instances: []eureka.InstanceInfo{
					{
						HostName: "http@A",
						App:      "A",
					},
				},
				proto: "grpc",
			},
			want: eureka.InstanceInfo{},
		},
		{
			name: "two",
			args: args{
				instances: []eureka.InstanceInfo{
					{
						HostName: "http@A",
						App:      "A",
					},
					{
						HostName: "grpc@B",
						App:      "B",
					},
				},

				proto: "grpc",
			},
			want: eureka.InstanceInfo{
				HostName: "grpc@B",
				App:      "B",
			},
		},
		{
			name: "two-  both grpc",
			args: args{
				instances: []eureka.InstanceInfo{
					{
						HostName: "grpc@A",
						App:      "A",
					},
					{
						HostName: "grpc@B",
						App:      "B",
					},
				},

				proto: "grpc",
			},
			want: eureka.InstanceInfo{
				HostName: "grpc@A",
				App:      "A",
			},
		},
		{
			name: "two - empty proto",
			args: args{
				instances: []eureka.InstanceInfo{
					{
						HostName: "http@A",
						App:      "A",
					},
					{
						HostName: "grpc@B",
						App:      "B",
					},
				},

				proto: "",
			},
			want: eureka.InstanceInfo{
				HostName: "http@A",
				App:      "A",
			},
		},
		{
			name: "empty",
			args: args{
				instances: []eureka.InstanceInfo{},
				proto:     "grpc",
			},
			want: eureka.InstanceInfo{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := selectOneOf(tt.args.instances, tt.args.proto); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("selectOneOf() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestResolveApplication(t *testing.T) {
	type args struct {
		fqServiceName    string
		proto            string
		includeWormholes bool
	}
	tests := []struct {
		name    string
		args    args
		want    eureka.InstanceInfo
		wantErr bool
	}{
		{
			name: "one correct",
			args: args{
				fqServiceName:    "DE.TU-BERLIN.COUNT",
				proto:            "grpc",
				includeWormholes: false,
			},
			want: eureka.InstanceInfo{
				HostName: "grpc@Theofaniss-iMac.fritz.box",
				App:      "DE.TU-BERLIN.COUNT",
			},
			wantErr: false,
		},
		{
			name: "one correct",
			args: args{
				fqServiceName:    "DE.TU-BERLIN.ERR",
				proto:            "grpc",
				includeWormholes: false,
			},
			want: eureka.InstanceInfo{
				HostName: "grpc@localhost",
				App:      "DE.TU-BERLIN.WH",
			},
			wantErr: false,
		},
		{
			name: "unknown",
			args: args{
				fqServiceName:    "PDF2TEXT",
				proto:            "grpc",
				includeWormholes: false,
			},
			want:    eureka.InstanceInfo{},
			wantErr: true,
		},
	}

	resolver := eureka.NewClient([]string{"http://eureka:8761/eureka"})
	cl := resolver.GetHttpClient()
	httpmock.ActivateNonDefault(cl)
	defer httpmock.DeactivateAndReset()

	COUNT, err := ioutil.ReadFile("../test/http/eureka/APP_DE_TU-BERLIN_COUNT.xml")
	if err != nil {
		t.Fatal(err)
	}
	responder := httpmock.NewBytesResponder(200, COUNT)
	httpmock.RegisterResponder("GET", "http://eureka:8761/eureka/apps/DE.TU-BERLIN.COUNT", responder)

	ALL, err := ioutil.ReadFile("../test/http/eureka/APP_ALL.xml")
	if err != nil {
		t.Fatal(err)
	}
	responder = httpmock.NewBytesResponder(200, ALL)
	httpmock.RegisterResponder("GET", "http://eureka:8761/eureka/apps", responder)

	ERROR, err := ioutil.ReadFile("../test/http/eureka/APP_ERR.xml")
	if err != nil {
		t.Fatal(err)
	}
	responder = httpmock.NewBytesResponder(404, ERROR)
	httpmock.RegisterResponder("GET", "http://eureka:8761/eureka/apps/DE.TU-BERLIN.ERR", responder)

	WH, err := ioutil.ReadFile("../test/http/eureka/APP_DE_TU-BERLIN_WH.xml")
	if err != nil {
		t.Fatal(err)
	}
	responder = httpmock.NewBytesResponder(200, WH)
	httpmock.RegisterResponder("GET", "http://eureka:8761/eureka/apps/DE.TU-BERLIN.WH", responder)

	httpmock.RegisterResponder("GET", `=~^http://eureka:8761/eureka/apps/.+\z`,
		httpmock.NewStringResponder(200, `{"id": 1, "name": "My Great Article"}`))

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			got, err := ResolveApplication(resolver, tt.args.fqServiceName, tt.args.proto, tt.args.includeWormholes)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResolveApplication() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !shallowEqualsInstanceID(got, tt.want) {
				t.Errorf("ResolveApplication() = %v, want %v", got, tt.want)
			}
		})
	}
}

func shallowEqualsInstanceID(got, want eureka.InstanceInfo) bool {
	if got.HostName != want.HostName || got.App != want.App {
		return false
	}
	return true
}

func TestBuildFQWormhole(t *testing.T) {
	type args struct {
		fqName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "normal",
			args: args{
				fqName: "DE.TU-BERLIN.COUNT",
			},
			want: "DE.TU-BERLIN.WH",
		},
		{
			name: "short",
			args: args{
				fqName: "DE.COUNT",
			},
			want: "DE.WH",
		},
		{
			name: "No scope",
			args: args{
				fqName: "COUNT",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildFQWormhole(tt.args.fqName); got != tt.want {
				t.Errorf("BuildFQWormhole() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWormholeResolveApplication(t *testing.T) {
	type args struct {
		scope            string
		fqServiceName    string
		proto            string
		includeWormholes bool
	}
	tests := []struct {
		name    string
		args    args
		want    eureka.InstanceInfo
		wantErr bool
	}{
		{
			name: "accessible via wh",
			args: args{
				scope:            "DE.TU-BERLIN",
				fqServiceName:    "DE.TU-BERLIN.ECHO",
				proto:            "grpc",
				includeWormholes: false,
			},
			want: eureka.InstanceInfo{
				HostName: "grpc@localhost",
				App:      "DE.TU-BERLIN.WH",
			},
			wantErr: false,
		},
		{
			name: "direct",
			args: args{
				scope:            "",
				fqServiceName:    "HTML2TEXT",
				proto:            "grpc",
				includeWormholes: false,
			},
			want: eureka.InstanceInfo{
				HostName: "grpc@localhost",
				App:      "HTML2TEXT",
			},
			wantErr: false,
		},
		{
			name: "not existent",
			args: args{
				scope:            "",
				fqServiceName:    "ERR",
				proto:            "grpc",
				includeWormholes: false,
			},
			want:    eureka.InstanceInfo{},
			wantErr: true,
		},
	}
	resolver := eureka.NewClient([]string{"http://eureka:8761/eureka"})
	cl := resolver.GetHttpClient()
	httpmock.ActivateNonDefault(cl)
	defer httpmock.DeactivateAndReset()

	COUNT, err := ioutil.ReadFile("../test/http/eureka/APP_DE_TU-BERLIN_COUNT.xml")
	if err != nil {
		t.Fatal(err)
	}
	responder := httpmock.NewBytesResponder(200, COUNT)
	httpmock.RegisterResponder("GET", "http://eureka:8761/eureka/apps/DE.TU-BERLIN.COUNT", responder)

	HTML2TEXT, err := ioutil.ReadFile("../test/http/eureka/APP_HTML2TEXT.xml")
	if err != nil {
		t.Fatal(err)
	}
	responder = httpmock.NewBytesResponder(200, HTML2TEXT)
	httpmock.RegisterResponder("GET", "http://eureka:8761/eureka/apps/HTML2TEXT", responder)

	ALL, err := ioutil.ReadFile("../test/http/eureka/APP_ALL.xml")
	if err != nil {
		t.Fatal(err)
	}
	responder = httpmock.NewBytesResponder(200, ALL)
	httpmock.RegisterResponder("GET", "http://eureka:8761/eureka/apps", responder)

	ERROR, err := ioutil.ReadFile("../test/http/eureka/APP_ERR.xml")
	if err != nil {
		t.Fatal(err)
	}
	responder = httpmock.NewBytesResponder(404, ERROR)
	httpmock.RegisterResponder("GET", "http://eureka:8761/eureka/apps/DE.TU-BERLIN.ERR", responder)
	httpmock.RegisterResponder("GET", "http://eureka:8761/eureka/apps/ERR", responder)

	WH, err := ioutil.ReadFile("../test/http/eureka/APP_DE_TU-BERLIN_WH.xml")
	if err != nil {
		t.Fatal(err)
	}
	responder = httpmock.NewBytesResponder(200, WH)
	httpmock.RegisterResponder("GET", "http://eureka:8761/eureka/apps/DE.TU-BERLIN.WH", responder)

	httpmock.RegisterResponder("GET", `=~^http://eureka:8761/eureka/apps/.+\z`,
		httpmock.NewStringResponder(200, `{"id": 1, "name": "My Great Article"}`))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := WormholeResolveApplication(resolver, tt.args.scope, tt.args.fqServiceName, tt.args.proto, tt.args.includeWormholes)
			if (err != nil) != tt.wantErr {
				t.Errorf("WormholeResolveApplication() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !shallowEqualsInstanceID(got, tt.want) {
				t.Errorf("WormholeResolveApplication() = %v, want %v", got, tt.want)
			}
		})
	}
}
