// +build api_examples

package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
)

func init() {
	_httpPort = startServer()
}

func TestGetCalls(t *testing.T) {

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
			t.Logf("%v", url)
			t.Logf(" --> %v", resp.StatusCode)
			defer resp.Body.Close()

			body, _ := ioutil.ReadAll(resp.Body)
			t.Logf("%v", string(body))
		})
	}
}

func TestPostCalls(t *testing.T) {
	tests := []struct {
		name     string
		endpoint string
		postBody DocRequest
	}{
		{
			"transform document",
			"/v1/document/transform",
			DocRequest{
				"abc", []byte(`2+3+4+5+(-1)`), "count", nil,
			},
		},
		{
			"transform document short",
			"/v1/document/transform",
			DocRequest{
				"abc", []byte(`2+3`), "add", nil,
			},
		},
	}

	//	_httpPort := startServer()

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			url := baseURL + ":" + strconv.Itoa(_httpPort) + tt.endpoint
			b, _ := json.Marshal(tt.postBody)
			resp, _ := http.Post(url, "application/json", bytes.NewBuffer(b))
			t.Logf("%v", url)
			t.Logf("Content-Type: %v", "application/json")
			t.Logf("%v", string(b))
			t.Logf(" --> %v", resp.StatusCode)
			defer resp.Body.Close()

			body, _ := ioutil.ReadAll(resp.Body)
			t.Logf("%v", string(body))
		})
	}
}

func TestPostTransformPipe(t *testing.T) {

	tests := []struct {
		name     string
		endpoint string
		postBody []DocRequest
	}{
		{
			"transform document",
			"/v1/document/transform-pipe",
			[]DocRequest{{
				"abc", []byte(`Hello World
				`), "count", nil,
			}},
		},
	}

	_httpPort := startServer()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := baseURL + ":" + strconv.Itoa(_httpPort) + tt.endpoint
			b, _ := json.Marshal(tt.postBody)
			resp, _ := http.Post(url, "application/json", bytes.NewBuffer(b))
			t.Logf("%v", url)
			t.Logf("Content-Type: %v", "application/json")
			t.Logf("%v", string(b))
			t.Logf(" --> %v", resp.StatusCode)
			defer resp.Body.Close()

			body, _ := ioutil.ReadAll(resp.Body)
			t.Logf("%v", string(body))
		})
	}
}
