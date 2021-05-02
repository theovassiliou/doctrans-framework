package main

import (
	"strings"
	"time"

	"github.com/jpillora/opts"
	"github.com/mitchellh/go-homedir"
	"github.com/theovassiliou/go-eureka-client/eureka"

	log "github.com/sirupsen/logrus"
	"github.com/theovassiliou/doctrans-framework/dtaservice"
	dta "github.com/theovassiliou/doctrans-framework/dtaservice"
	"github.com/theovassiliou/doctrans-framework/instanceid"
	aux "github.com/theovassiliou/doctrans-framework/ipaux"
	wh "github.com/theovassiliou/doctrans-framework/services/wormhole/whimplementation"
)

//set this via ldflags (see https://stackoverflow.com/q/11354518)
// version is the current version number as tagged via git tag 1.0.0 -m 'A message'
var (
	version   = dta.Version
	commit    string
	branch    string
	cmdName   = "wormhole"
	startTime time.Time
)

func init() {
	startTime = time.Now()
}

const (
	appName = "DE.TU-BERLIN.WH"
	dtaType = "Gateway"
	scope   = "DE.TU-BERLIN"
)

var resolver *eureka.Client

type whCmdLineOptions struct {
	dta.DocTransServerOptions
	dta.DocTransServerGenericOptions
	ResolverURL          string `opts:"group=Resolver" help:"Resolver URL"`
	ResolverRegistration bool   `opts:"group=Resolver" help:"Register in addition also to the resolver"`
	Scope                string `opts:"group=Wormhole" help:"Defines the namespace of the galaxy this wormhole is service"`
	appName              string
}

func main() {
	workingHomeDir, _ := homedir.Dir()
	homepageURL := "https://github.com/theovassiliou/doctrans/blob/master/wormhole/README.md"
	gwOptions := whCmdLineOptions{}
	gwOptions.CfgFile = workingHomeDir + "/.dta/" + appName + "/config.json"
	gwOptions.LogLevel = log.WarnLevel
	gwOptions.RegHostName = aux.GetHostname()
	gwOptions.RegistrarURL = "http://eureka:8761/eureka"
	gwOptions.ResolverURL = "http://eureka:8761/eureka"
	gwOptions.Scope = scope

	opts.New(&gwOptions).
		Repo("github.com/theovassiliou/doctrans").
		ConfigPath(gwOptions.CfgFile).
		Version(dtaservice.FormatFullVersion(cmdName, version, branch, commit)).
		Parse()

	if gwOptions.LogLevel != 0 {
		log.SetLevel(gwOptions.LogLevel)
	}

	if gwOptions.Scope != scope {
		gwOptions.appName = strings.TrimSuffix(gwOptions.Scope, ".") + ".WH"
	}

	if gwOptions.appName != "" {
		gwOptions.CfgFile = workingHomeDir + "/.dta/" + gwOptions.appName + "/config.json"
		opts.New(&gwOptions).
			Repo("github.com/theovassiliou/doctrans").
			ConfigPath(gwOptions.CfgFile).
			Version(dtaservice.FormatFullVersion(cmdName, version, branch, commit)).
			Parse()
	} else {
		gwOptions.appName = appName
	}

	var _grpcGateway, _httpGateway dta.IDocTransServer
	// Calc Configuration
	registerGRPC, registerHTTP := determineServerConfig(gwOptions)
	if registerGRPC {
		_grpcGateway = newWormholeService(gwOptions, gwOptions.Scope, gwOptions.appName, "grpc")
	}
	if registerHTTP {
		_httpGateway = newWormholeService(gwOptions, gwOptions.Scope, gwOptions.appName, "http")
	}

	// create client resolver
	gg := _grpcGateway.GetDocTransServer()
	gg.SetResolver(eureka.NewClient([]string{
		gwOptions.ResolverURL,
	}))
	dta.LaunchServices(_grpcGateway, _httpGateway, gwOptions.appName, dtaType, homepageURL, gwOptions.DocTransServerOptions)
}

func newWormholeService(options whCmdLineOptions, scope, appName, proto string) dta.IDocTransServer {
	gw := wh.Wormhole{
		UnimplementedDTAServerServer: dta.UnimplementedDTAServerServer{},
		GenDocTransServer:            dta.GenDocTransServer{AppName: appName, DtaType: dtaType, Proto: proto, XInstanceIDprefix: instanceid.BuildVBC(appName, version, branch, commit)},
		IDocTransServer:              nil,
		Scope:                        scope,
	}
	gw.AppName = appName
	if !options.XInstanceID {
		gw.XInstanceIDprefix = instanceid.BuildVBC(appName, version, branch, commit)
		gw.XInstanceIDstartTime = startTime
	}
	return &gw
}

func determineServerConfig(gwOptions whCmdLineOptions) (registerGRPC, registerHTTP bool) {
	if (!gwOptions.HTTP && !gwOptions.GRPC) || gwOptions.GRPC {
		registerGRPC = true
	}

	if (!gwOptions.HTTP && !gwOptions.GRPC) || gwOptions.HTTP {
		registerHTTP = true
	}
	return
}
