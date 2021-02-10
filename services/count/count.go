package main

// A simple implemenation of using the Golang DocTrans Framework
import (
	"fmt"
	"time"

	"github.com/jpillora/opts"
	homedir "github.com/mitchellh/go-homedir"

	log "github.com/sirupsen/logrus"
	"github.com/theovassiliou/doctrans-framework/dtaservice"
	dta "github.com/theovassiliou/doctrans-framework/dtaservice"
	aux "github.com/theovassiliou/doctrans-framework/ipaux"
	service "github.com/theovassiliou/doctrans-framework/services/count/serviceimplementation"
)

//set this via ldflags (see https://stackoverflow.com/q/11354518)
// version is the current version number as tagged via git tag 1.0.0 -m 'A message'
var (
	version   = dta.Version
	commit    string
	branch    string
	cmdName   = "count"
	startTime time.Time
)

func init() {
	startTime = time.Now()
}

const (
	appName = "DE.TU-BERLIN.COUNT"
	dtaType = "Service"
)

type serviceCmdLineOptions struct {
	dta.DocTransServerOptions
	dta.DocTransServerGenericOptions
	LocalExecution string `opts:"group=Local Execution, short=x" help:"If set, execute the service locally once and read from this file"`
}

func main() {
	workingHomeDir, _ := homedir.Dir()
	homepageURL := dta.RepoName

	serviceOptions := serviceCmdLineOptions{}
	serviceOptions.CfgFile = workingHomeDir + "/.dta/" + appName + "/config.json"
	serviceOptions.Port = 50000
	serviceOptions.LogLevel = log.WarnLevel
	serviceOptions.HostName = aux.GetHostname()
	serviceOptions.RegistrarURL = "http://eureka:8761/eureka"

	opts.New(&serviceOptions).
		Repo("github.com/theovassiliou/doctrans").
		ConfigPath(serviceOptions.CfgFile).
		Version(dtaservice.FormatFullVersion(cmdName, version, branch, commit)).
		Parse()

	if serviceOptions.LogLevel != 0 {
		log.SetLevel(serviceOptions.LogLevel)
	}

	if serviceOptions.LocalExecution != "" {
		s := service.DtaService{}
		s.AppName = appName
		transDoc := service.ExecuteWorkerLocally(s, serviceOptions.LocalExecution)
		fmt.Println(transDoc)
		return
	}

	var _grpcGateway, _httpGateway dta.IDocTransServer
	// Calc Configuration
	registerGRPC, registerHTTP := determineServerConfig(serviceOptions)
	if registerGRPC {
		_grpcGateway = newDtaService(serviceOptions, appName, "grpc")
	}
	if registerHTTP {
		_httpGateway = newDtaService(serviceOptions, appName, "http")
	}

	dta.LaunchServices(_grpcGateway, _httpGateway, appName, dtaType, homepageURL, serviceOptions.DocTransServerOptions)
}
func newDtaService(options serviceCmdLineOptions, appName, proto string) dta.IDocTransServer {
	gw := service.DtaService{
		GenDocTransServer: dta.GenDocTransServer{
			AppName:           appName,
			DtaType:           dtaType,
			Proto:             proto,
			XInstanceIDprefix: buildXIIDprefix(),
		},
	}
	gw.AppName = appName
	if !options.XInstanceID {
		gw.XInstanceIDprefix = buildXIIDprefix()
		gw.XInstanceIDstartTime = startTime
	}
	return &gw
}

func determineServerConfig(gwOptions serviceCmdLineOptions) (registerGRPC, registerHTTP bool) {
	if (!gwOptions.HTTP && !gwOptions.GRPC) || gwOptions.GRPC {
		registerGRPC = true
	}

	if (!gwOptions.HTTP && !gwOptions.GRPC) || gwOptions.HTTP {
		registerHTTP = true
	}
	return
}

func buildXIIDprefix() string {
	return appName + "/" + version + "/" + branch + "-" + commit + "/"
}
