package dtaservice

import (
	"context"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
)

func Test_calcStatusURL(t *testing.T) {
	type args struct {
		instanceID string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"simple",
			args{"x"},
			"http://x/status",
		},
		{
			"hostname",
			args{"www.test.com"},
			"http://www.test.com/status",
		},
		{
			"ipAddress",
			args{"127.0.0.1"},
			"http://127.0.0.1/status",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calcStatusURL(tt.args.instanceID); got != tt.want {
				t.Errorf("calcStatusURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_calcHealthURL(t *testing.T) {
	type args struct {
		instanceID string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"simple",
			args{"x"},
			"http://x/health",
		},
		{
			"hostname",
			args{"www.test.com"},
			"http://www.test.com/health",
		},
		{
			"ipAddress",
			args{"127.0.0.1"},
			"http://127.0.0.1/health",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calcHealthURL(tt.args.instanceID); got != tt.want {
				t.Errorf("calcHealthURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMuxHTTPIntoGrpc(t *testing.T) {

	// Just checking the default functionality

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	_httpListener, _httpPort := CreateListener(50000, 20)

	go MuxHTTPIntoGrpc(ctx, _httpListener, 0)

	// /status
	resp, err := http.Get("http://127.0.0.1:" + strconv.Itoa(_httpPort) + "/status")
	if err != nil {
		t.Errorf("err = %v, want %v", err, nil)
	}
	if resp.StatusCode != 200 {
		t.Errorf("StatusCode = %v, want %v", resp.StatusCode, 200)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if string(body) != "The service is alive" {
		t.Errorf("Body = %v, want %v", string(body), nil)
	}

	// /health
	resp, err = http.Get("http://127.0.0.1:" + strconv.Itoa(_httpPort) + "/health")
	if err != nil {
		t.Errorf("err = %v, want %v", err, nil)
	}
	if resp.StatusCode != 200 {
		t.Errorf("StatusCode = %v, want %v", resp.StatusCode, 200)
	}

	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)

	if string(body) != "The service is healthy as it is responding." {
		t.Errorf("Body = %v, want %v", string(body), nil)
	}

}
