package aux

import (
	"net"
	"testing"
)

func TestExternalIP(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping ExternalIP() in short mode")
	}
	got, err := ExternalIP()

	if err != nil {
		t.Errorf("ExternalIP() error = %v, wantErr %v", err, nil)
		return
	}

	ip := net.ParseIP(got)
	if ip.String() != got {
		t.Errorf("ExternalIP() != parsed IP (%v)", ip)
	}

}

func TestGetHostname(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping GetHostname() in short mode")
	}
	got := GetHostname()
	if got == "" {
		t.Errorf("GetHostname returne empty string. Please check.")
	}
}

func TestGetIPAdress(t *testing.T) {
	eIP, _ := ExternalIP()
	got := GetIPAdress()

	if got != eIP {
		t.Errorf("GetIPAdress() != ExternalIP() %v!=%v", got, eIP)
	}
}

func TestPublicIP(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping ExternalIP() in short mode")
	}

	public, err := PublicIP()

	if err != nil {
		t.Errorf("PublicIP() error = %v, wantErr %v", err, nil)
		return
	}

	external, err := ExternalIP()

	if public == external {
		t.Errorf("PublicIP(%v) == ExternalIP(%v). Check!", public, external)
	}
}
