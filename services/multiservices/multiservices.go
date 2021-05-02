package main

// A simple implemenation of using the Golang DocTrans Framework
import (
	"fmt"
	"sync"
	"time"

	"github.com/jpillora/opts"
	homedir "github.com/mitchellh/go-homedir"

	log "github.com/sirupsen/logrus"
	"github.com/theovassiliou/doctrans-framework/dtaservice"
	dta "github.com/theovassiliou/doctrans-framework/dtaservice"
	aux "github.com/theovassiliou/doctrans-framework/ipaux"
	echo "github.com/theovassiliou/doctrans-framework/services/multiservices/echoimplementation"
	h2t "github.com/theovassiliou/doctrans-framework/services/multiservices/html2textimplementation"
)

//set this via ldflags (see https://stackoverflow.com/q/11354518)
// version is the current version number as tagged via git tag 1.0.0 -m 'A message'
var (
	version   = dta.Version
	commit    string
	branch    string
	cmdName   = "multiservice"
	startTime time.Time
)

func init() {
	startTime = time.Now()
}

type service struct {
	serviceName    string
	serviceCreator func(options serviceCmdLineOptions, appName, proto string) dta.IDocTransServer
}

var appName = "MULTISERVICE"
var localServices = []service{
	{
		"ECHO",
		newEchoService,
	},
	{
		"HTML2TEXT",
		newHTML2TextService,
	},
}

const (
	dtaType = "Service"
)

type serviceCmdLineOptions struct {
	dta.DocTransServerOptions
	dta.DocTransServerGenericOptions
	LocalExecution string `opts:"group=Local Execution, short=x" help:"If set, execute the service locally once and read from this file"`
	HTML2Text      bool   `opts:"group=Local Execution, short=1" help:"If set, use HTML2TEXT service"`
	Echo           bool   `opts:"group=Local Execution, short=2" help:"If set, use ECHO service"`
}

func main() {
	workingHomeDir, _ := homedir.Dir()
	homepageURL := dta.RepoName

	serviceOptions := serviceCmdLineOptions{}
	serviceOptions.CfgFile = workingHomeDir + "/.dta/" + cmdName + "/config.json"
	serviceOptions.Port = 50000
	serviceOptions.LogLevel = log.WarnLevel
	serviceOptions.RegHostName = aux.GetHostname()
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
		if serviceOptions.HTML2Text {
			s := h2t.DtaService{}
			s.AppName = appName
			transDoc := h2t.ExecuteWorkerLocally(s, serviceOptions.LocalExecution)
			fmt.Println(transDoc)
		}

		if serviceOptions.Echo {
			s := echo.DtaService{}
			s.AppName = appName
			transDoc := echo.ExecuteWorkerLocally(s, serviceOptions.LocalExecution)
			fmt.Println(transDoc)
		}
		return
	}

	var wg sync.WaitGroup

	for _, s := range localServices {
		var _grpcGateway, _httpGateway dta.IDocTransServer
		// Calc Configuration
		registerGRPC, registerHTTP := determineServerConfig(serviceOptions)
		if registerGRPC {
			_grpcGateway = s.serviceCreator(serviceOptions, s.serviceName, "grpc")
		}
		if registerHTTP {
			_httpGateway = s.serviceCreator(serviceOptions, s.serviceName, "http")
		}
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			dta.LaunchServices(_grpcGateway, _httpGateway, s.serviceName, dtaType, homepageURL, serviceOptions.DocTransServerOptions)
		}(&wg)
		wg.Add(1)
		time.Sleep(time.Millisecond * 500)
	}

	wg.Wait()
}

func newHTML2TextService(options serviceCmdLineOptions, appName, proto string) dta.IDocTransServer {
	gw := h2t.DtaService{
		GenDocTransServer: dta.GenDocTransServer{
			AppName:           appName,
			DtaType:           dtaType,
			Proto:             proto,
			XInstanceIDprefix: buildXIIDprefix(appName),
		},
	}
	gw.AppName = appName
	if !options.XInstanceID {
		gw.XInstanceIDprefix = buildXIIDprefix(appName)
		gw.XInstanceIDstartTime = startTime
	}
	return &gw
}

func newEchoService(options serviceCmdLineOptions, appName, proto string) dta.IDocTransServer {
	gw := echo.DtaService{
		GenDocTransServer: dta.GenDocTransServer{
			AppName:           appName,
			DtaType:           dtaType,
			Proto:             proto,
			XInstanceIDprefix: buildXIIDprefix(appName),
		},
	}
	gw.AppName = appName
	if !options.XInstanceID {
		gw.XInstanceIDprefix = buildXIIDprefix(appName)
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

func buildXIIDprefix(appName string) string {
	return appName + "/" + version + "/" + branch + "-" + commit + "/"
}
