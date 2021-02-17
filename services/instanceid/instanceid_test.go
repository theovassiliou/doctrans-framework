package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"testing"

	_struct "github.com/golang/protobuf/ptypes/struct"
	log "github.com/sirupsen/logrus"

	dta "github.com/theovassiliou/doctrans-framework/dtaservice"
)

type DocRequest struct {
	FileName    string          `json:"file_name,omitempty"`
	Document    []byte          `json:"document,omitempty"`
	ServiceName string          `json:"service_name,omitempty"`
	Options     *_struct.Struct `json:"options,omitempty"`
}

type DocResponse struct {
	Document    []byte   `json:"document,omitempty"`
	TransOutput []string `json:"trans_output,omitempty"`
	Error       []string `json:"error,omitempty"`
}

const baseURL = "http://127.0.0.1"

var _httpPort int

func init() {
	_httpPort = startServer()
}

func TestBasicStatus(t *testing.T) {

	tests := []struct {
		name       string
		endpoint   string
		want       string
		statusCode int
	}{
		{
			"status",
			"/status",
			"The service is alive",
			200,
		},
		{
			"health",
			"/health",
			"The service is healthy as it is responding.",
			200,
		},
	}

	//	_httpPort := startServer()

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			url := baseURL + ":" + strconv.Itoa(_httpPort) + tt.endpoint
			resp, _ := http.Get(url)
			if resp.StatusCode != tt.statusCode {
				t.Errorf("%v() = %v, want %v", tt.endpoint, resp.StatusCode, tt.statusCode)
			}
			defer resp.Body.Close()

			body, _ := ioutil.ReadAll(resp.Body)

			if string(body) != tt.want {
				t.Errorf("%v() = %v, want %v", tt.endpoint, string(body), tt.want)
			}
		})
	}
}

func TestCountDtaGet(t *testing.T) {

	tests := []struct {
		name       string
		endpoint   string
		want       string
		statusCode int
	}{
		{
			"list",
			"/v1/service/list",
			`{"services":["INSTANCEID"]}`,
			200,
		},
		{
			"options",
			"/v1/service/options",
			`{"error":"method Options not implemented","code":12,"message":"method Options not implemented"}`,
			501,
		},
	}

	_httpPort := startServer()

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			url := baseURL + ":" + strconv.Itoa(_httpPort) + tt.endpoint
			resp, _ := http.Get(url)
			if resp.StatusCode != tt.statusCode {
				t.Errorf("%v() = %v, want %v", tt.endpoint, resp.StatusCode, tt.statusCode)
			}
			defer resp.Body.Close()

			body, _ := ioutil.ReadAll(resp.Body)

			if string(body) != tt.want {
				t.Errorf("%v() = %v, want %v", tt.endpoint, string(body), tt.want)
			}
		})
	}
}

func TestCountDtaTransform(t *testing.T) {

	tests := []struct {
		name       string
		endpoint   string
		postBody   DocRequest
		want       DocResponse
		statusCode int
	}{
		{
			"transform document",
			"/v1/document/transform",
			DocRequest{
				"abc", []byte(`MsA/1.1/xxx%22s(msC/1.4%5555s+msD/2.2%23234s)`), "instance-id", nil,
			},
			DocResponse{
				[]byte(`.
└── [22s]  MsA/1.1
    ├── [5555s]  msC/1.4
    └── [23234s]  msD/2.2
`), nil, nil,
			},
			200,
		},
	}

	//	_httpPort := startServer()

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			url := baseURL + ":" + strconv.Itoa(_httpPort) + tt.endpoint
			b, _ := json.Marshal(tt.postBody)
			resp, _ := http.Post(url, "application/json", bytes.NewBuffer(b))
			if resp.StatusCode != tt.statusCode {
				t.Errorf("%v() = %v, want %v", tt.endpoint, resp.StatusCode, tt.statusCode)
			}
			defer resp.Body.Close()

			body, _ := ioutil.ReadAll(resp.Body)

			var rb DocResponse
			json.Unmarshal(body, &rb)
			if !reflect.DeepEqual(rb, tt.want) {
				t.Errorf("%v() = %v, want %v", tt.endpoint, rb, tt.want)

			}
		})
	}
}

func TestCountDtaTransformPipe(t *testing.T) {

	tests := []struct {
		name       string
		endpoint   string
		postBody   []DocRequest
		want       DocResponse
		statusCode int
	}{
		{
			"transform document",
			"/v1/document/transform-pipe",
			[]DocRequest{{
				"abc", []byte(`Hello World
				`), "count", nil,
			}},
			DocResponse{},
			400,
		},
	}

	_httpPort := startServer()

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			url := baseURL + ":" + strconv.Itoa(_httpPort) + tt.endpoint
			b, _ := json.Marshal(tt.postBody)
			resp, _ := http.Post(url, "application/json", bytes.NewBuffer(b))
			if resp.StatusCode != tt.statusCode {
				t.Errorf("%v() = %v, want %v", tt.endpoint, resp.StatusCode, tt.statusCode)
			}
			defer resp.Body.Close()

			body, _ := ioutil.ReadAll(resp.Body)

			var rb DocResponse
			json.Unmarshal(body, &rb)
			if !reflect.DeepEqual(rb, tt.want) {
				t.Errorf("%v() = %v, want %v", tt.endpoint, rb, tt.want)

			}
		})
	}
}

func startServer() (httpPort int) {
	log.SetLevel(log.ErrorLevel)

	ctx := context.Background()

	// Uncommented this to be package testable
	// ctx, _ = context.WithCancel(ctx)
	// defer cancel()

	_grpcGateway := newDtaService(serviceCmdLineOptions{}, appName, "grpc")
	_ = newDtaService(serviceCmdLineOptions{}, appName, "http")
	_grpcListener, _grpcPort := dta.CreateListener(50000, 20)
	_httpListener, _httpPort := dta.CreateListener(_grpcPort, 20)

	go dta.StartGrpcServer(_grpcListener, _grpcGateway)
	go dta.MuxHTTPIntoGrpc(ctx, _httpListener, _grpcPort)
	return _httpPort
}
