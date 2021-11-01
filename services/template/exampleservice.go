package main

import (
	"time"

	"github.com/gin-gonic/gin"
)

const (
	VERSION        = "1.0"
	BRANCH         = "src"
	COMMIT         = ""
	SERVICENAME    = "exampleservice"
	THISSERVICE    = SERVICENAME + "/" + VERSION + "/-1s"
	CMDLINEVERSION = VERSION + " " + "(git: " + BRANCH + " " + COMMIT + ")"
	APIVERSION     = "/v1"
)

var startTime time.Time

func init() {
	startTime = time.Now()
}

func main() {
	r := gin.Default()
	apiv1 := r.Group(APIVERSION)

	document := apiv1.Group(("/document"))

	// /document/transform
	document.POST("/transform", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"document":     "example response",
			"trans_output": []string{"transformed successful"},
			"error":        []string{},
		})
	})

	// /document/transform-pipe
	document.POST("/transform-pipe", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"document":     "",
			"trans_output": []string{},
			"error":        []string{"not implemented"},
		})
	})

	service := apiv1.Group("/service")
	// /service/list
	service.GET("/list", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"services": []string{SERVICENAME},
		})
	})
	r.Run(":50002")
}
