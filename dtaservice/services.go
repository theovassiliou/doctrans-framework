package dtaservice

import (
	"net"
	"strconv"

	log "github.com/sirupsen/logrus"
)

// CreateListener creates a net.listener on port that is yet to be defined. startPort contains the port number where listen should start with and
// continues until port number  startPort+maxPortSeek to open a listener. If it finally fails nil is returned as listener and the last portnumber tries as
// return port
func CreateListener(startPort, maxPortSeek int) (net.Listener, int) {
	var lis net.Listener
	var err error
	initialPort := startPort

	for i := 0; i < maxPortSeek; i++ {
		log.WithFields(log.Fields{"Service": "Server", "Status": "Trying"}).Infof("Trying to listen on port %d", (startPort + i))
		lis, err = net.Listen("tcp", ":"+strconv.Itoa(startPort+i))
		if err == nil {
			startPort += i
			log.WithFields(log.Fields{"Service": "Server", "Status": "Listening"}).Infof("Using port %d to listen for dta", startPort)
			i = maxPortSeek
		}
	}

	if err != nil {
		log.WithFields(log.Fields{"Service": "Server", "Status": "Abort"}).Infof("Failed to finally open ports between %d and %d", startPort, startPort+maxPortSeek)
		log.WithFields(log.Fields{"Service": "Server", "Status": "Abort"}).Fatalf("failed to listen: %v", err)
		return nil, startPort
	}

	if startPort != initialPort {
		log.WithFields(log.Fields{"Service": "Server", "Status": "main"}).Warnf("Listing on port %v instead on configured, but used port %v\n", initialPort, startPort)
	}

	log.WithFields(log.Fields{"Service": "Server", "Status": "main"}).Debugln("Opend successfull a port")
	return lis, startPort
}

// ApplicationName returns the application name of the service
func (dtas *GenDocTransServer) ApplicationName() string {
	return dtas.AppName
}

func calcHostName(proto, hostName string) string {
	return proto + "@" + hostName
}
