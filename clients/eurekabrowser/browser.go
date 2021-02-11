package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jpillora/opts"
	"github.com/theovassiliou/doctrans-framework/dtaservice"
	dta "github.com/theovassiliou/doctrans-framework/dtaservice"

	log "github.com/sirupsen/logrus"
	"github.com/theovassiliou/go-eureka-client/eureka"
	"github.com/xlab/treeprint"
)

// REPONAME is the name of the repo at github
var REPONAME = "doctrans-framework"
var (
	version   = dta.Version
	commit    string
	branch    string
	cmdName   = "eurekabrowser"
	startTime time.Time
)

// REPO is the full github name
var REPO = "github.com/theovasiliou/" + REPONAME

type config struct {
	EurekaURL []string  `type:"arg" min:"0" name:"file" help:"eureka servers to query, http://eureka:8761/eureka if none"`
	LogLevel  log.Level `help:"Log level, one of panic, fatal, error, warn or warning, info, debug, trace"`
}

var conf = config{}

func main() {
	conf = config{
		LogLevel: log.InfoLevel,
	}

	//parse config
	opts.New(&conf).
		Repo(REPO).
		Version(dtaservice.FormatFullVersion(cmdName, version, branch, commit)).
		Parse()

	if len(conf.EurekaURL) == 0 {
		conf.EurekaURL = []string{"http://eureka:8761/eureka"}
	}

	log.SetLevel(conf.LogLevel)
	base := collectTree{treeprint.New()}
	base.t.SetValue("all servers")
	for _, c := range conf.EurekaURL {

		client := eureka.NewClient([]string{c})
		applications, _ := client.GetApplications() // Retrieves all applications from eureka server(s)

		tree := collectTree{base.t.AddBranch(".")}
		tree.t.SetValue(strings.TrimRight(c, "/eureka"))

		for _, app := range applications.Applications {
			app.Accept(tree)
		}
	}
	fmt.Println((base.t).String())
}

type collectTree struct {
	t treeprint.Tree
}

func (t collectTree) VisitForApplication(a eureka.Application) {
	x := collectTree{t.t.AddBranch(a.Name)}
	for _, i := range a.Instances {
		i.Accept(x)
	}
}

func (t collectTree) VisitForInstance(i eureka.InstanceInfo) {
	metaInfo := strconv.Itoa(int(time.Since(time.Unix(int64(i.LastUpdatedTimestamp)/1000, 0)).Seconds())) + "s"
	branchName := i.IpAddr + ":" + i.Port.Port
	t.t.AddMetaBranch(metaInfo, branchName)
}
