/*
Copyright (c) 2019 Theofanis Vassiliou-Gioles

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

//go:generate protoc -I ../dtaservice --go_out=plugins=grpc:../dtaservice ../dtaservice/dtaservice.proto

// Package main implements a client for DtaService.
package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/jpillora/opts"

	log "github.com/sirupsen/logrus"
	"github.com/theovassiliou/doctrans-framework/dtaservice"
	dta "github.com/theovassiliou/doctrans-framework/dtaservice"
	pb "github.com/theovassiliou/doctrans-framework/dtaservice"

	"github.com/theovassiliou/go-eureka-client/eureka"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
)

// REPONAME is the name of the repo at github
var REPONAME = "doctrans-framework"
var (
	version   = dta.Version
	commit    string
	branch    string
	cmdName   = "grpcclient"
	startTime time.Time
)

// REPO is the full github name
var REPO = "github.com/theovasiliou/" + REPONAME

type config struct {
	FileName       []string  `type:"arg" min:"0" name:"file" help:"the file to be uploaded"`
	EurekaURL      string    `help:"if set the indicated eureka server will be used to find DTA-GW"`
	ServiceName    string    `help:"The service to be used"`
	ServiceAddress string    `help:"Address and port of the server to connect"`
	ListServices   bool      `help:"List all the services accessible"`
	LogLevel       log.Level `help:"Log level, one of panic, fatal, error, warn or warning, info, debug, trace"`
}

var conf = config{}

//set this via ldflags (see https://stackoverflow.com/q/11354518)
// version is the current version number as tagged via git tag 1.0.0 -m 'A message'

func check(e error) {
	if e != nil {
		panic(e)
	}
}

const dtaGwID = "BERLIN.VASSILIOU-POHL.GW"

func main() {
	conf = config{
		ServiceName: "DE.TU-BERLIN.COUNT",
		EurekaURL:   "http://eureka:8761/eureka",
		LogLevel:    log.InfoLevel,
	}

	//parse config
	opts.New(&conf).
		Repo(REPO).
		Version(dtaservice.FormatFullVersion(cmdName, version, branch, commit)).
		Parse()

	log.SetLevel(conf.LogLevel)

	// Set up a connection to the server.
	log.Infof("Requesting service %s", conf.ServiceName)

	// We have to identify the server to contact
	// We have to possibilities
	//  a) via registry (the normal case)
	//  b) direct, more for testing purposes

	// 	a) via resolver is assumed if no server is given
	//  - contact the well-known resolver
	if conf.ServiceAddress == "" {
		log.Infof("Will contact registry at %s\n", conf.EurekaURL)

		client := eureka.NewClient([]string{
			conf.EurekaURL, //From a spring boot based eureka server
			// add others servers here
		})

		// Let's find out whether we find the server serving this service.
		//  - ask for the service
		eService, err := client.GetApplication(conf.ServiceName)
		if err != nil || len(eService.Instances) == 0 {
			log.Infof("Could not find the service %s at eureka\n", conf.ServiceName)
		} else {
			conf.ServiceAddress = eService.Instances[1].IpAddr + ":" + eService.Instances[1].Port.Port
		}
		//  - if service is unknown ask for a gateway
		if conf.ServiceAddress == "" {
			log.Infof("Looking for a gateway %s \n", dtaGwID)

			gService, err := client.GetApplication(dtaGwID)
			if err != nil || len(gService.Instances) == 0 {
				log.Infof("Could not find a gateway %s at eureka\n", dtaGwID)
			} else {
				conf.ServiceAddress = gService.Instances[0].IpAddr + ":" + gService.Instances[0].Port.Port
				log.Infof("Found one at %s \n", conf.ServiceAddress)
			}
		}
	}
	log.Infof("Will contact %s for service for service %s\n", conf.ServiceAddress, conf.ServiceName)

	//  - contact identified server

	conn, err := grpc.Dial(conf.ServiceAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewDTAServerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if conf.ListServices {
		r, err := c.ListServices(ctx, &emptypb.Empty{})
		if err != nil {
			log.Fatalf("could not list services: %v", err)
		}

		fmt.Println(strings.Join(r.GetServices(), "\n"))
		os.Exit(0)
	}

	for _, fN := range conf.FileName {
		// Read content from file.
		doc, err := ioutil.ReadFile(fN)
		if err != nil {
			log.Warnln(err.Error())
			log.Warnln("Skipping", fN)
			continue
		}

		var header metadata.MD
		options := structpb.NewNullValue().GetStructValue()
		r, err := c.TransformDocument(ctx, &pb.TransformDocumentRequest{ServiceName: conf.ServiceName, FileName: fN, Document: doc, Options: options}, grpc.Header(&header))
		if err != nil {
			log.Fatalf("could not transform: %v", err)
		} else if r.GetError() != nil {
			fmt.Println(strings.Join(r.GetError(), "\n"))
			return
		}
		fmt.Println(fN)
		fmt.Println(string(r.GetDocument()))
		fmt.Printf("Received-Header: %#v\n", header)
	}
}