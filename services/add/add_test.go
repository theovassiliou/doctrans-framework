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

func TestAddDtaGet(t *testing.T) {

	tests := []struct {
		name       string
		endpoint   string
		want       string
		statusCode int
	}{
		{
			"list",
			"/v1/service/list",
			`{"services":["ADD"]}`,
			200,
		},
		{
			"options",
			"/v1/service/options",
			`{"error":"method Options not implemented","code":12,"message":"method Options not implemented"}`,
			501,
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
				"abc", []byte(`2+3`), "count", nil,
			},
			DocResponse{
				[]byte(`5`), nil, nil,
			},
			200,
		},
		{
			"adding minus",
			"/v1/document/transform",
			DocRequest{
				"abc", []byte(`2+(-3)`), "count", nil,
			},
			DocResponse{
				[]byte(`-1`), nil, nil,
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

func TestAddOperation(t *testing.T) {

	tests := []struct {
		name string

		ask        string
		result     string
		statusCode int
	}{
		{
			"simple",
			`2+3`,
			`5`,
			200,
		},
		{
			"simple with spaces",
			` 2+ 3 `,
			`5`,
			200,
		},
		{
			"one fixed number",
			` 2.0 + 3 `,
			`5`,
			200,
		},
		{
			"two fixed number",
			` 2.0 + 3.0 `,
			`5`,
			200,
		},
		{
			"two fixed number, fixed result",
			` 2.0 + 3.1 `,
			`5.1`,
			200,
		},
		{
			"two fixed number, float result",
			` 3423424232.0 + 3.123423423423 `,
			`3.4234242351234236e+09`,
			200,
		},
		{
			"two larger number float, float result",
			` 555555 +555555 `,
			`1.11111e+06`,
			200,
		},
		{
			"two larger number float, float result",
			` 5.00000000000001 +5 `,
			`10.00000000000001`,
			200,
		},
		{
			"not enought precision",
			` 5.000000000000001 - 5 `,
			`8.881784197001252e-16`,
			200,
		},
		{
			"multiple args",
			` 5 + 5 + 5 + 5 - 5 `,
			`15`,
			200,
		},
		// {
		// 	"multiple args",
		// 	` a + 5 + 5 + 5 - 5 `,
		// 	"",
		// 	200,
		// },
	}

	//	_httpPort := startServer()

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			url := baseURL + ":" + strconv.Itoa(_httpPort) + "/v1/document/transform"
			b, _ := json.Marshal(DocRequest{
				"abc", []byte(tt.ask), "add", nil,
			})
			resp, _ := http.Post(url, "application/json", bytes.NewBuffer(b))
			if resp.StatusCode != tt.statusCode {
				t.Errorf("%v() = %v, want %v", "/v1/document/transform", resp.StatusCode, tt.statusCode)
			}
			defer resp.Body.Close()

			body, _ := ioutil.ReadAll(resp.Body)
			var rb DocResponse
			json.Unmarshal(body, &rb)

			if !reflect.DeepEqual(rb, DocResponse{
				[]byte(tt.result), nil, nil,
			}) {
				t.Errorf("%v() = .%s., want .%v.", "/v1/document/transform", string(rb.Document), tt.result)

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
				"abc", []byte(`2+3`), "add", nil,
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
