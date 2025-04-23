// @author: yu-xiaoyao
// @github: https://github.com/yu-xiaoyao/jrebel-license-active-server
package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Config struct {
	Port           int
	OfflineDefault bool
	OfflineDays    int
	LogLevel       int64
	ExportSchema   string
	ExportHost     string
	NewIndex       bool
	BasePath       string
}

var config = &Config{
	Port:           12345,
	OfflineDefault: true,
	OfflineDays:    30,
	LogLevel:       Info,
	ExportSchema:   "http",
	ExportHost:     "", // default is request ip
	NewIndex:       true,
	BasePath:       "",
}

var logger = NewLogger(os.Stdout, Info, log.Ldate|log.Ltime)

func init() {
	portPtr := flag.Int("port", config.Port, "Server port, value range 1-65535")
	logLevelPtr := flag.Int64("logLevel", config.LogLevel, "Log level, value range 1-4")
	exportSchemaPtr := flag.String("exportSchema", config.ExportSchema, "Index page export schema (http or https)")
	exportHostPtr := flag.String("exportHost", "", "Index page export host, ip or domain")
	newIndexPtr := flag.Bool("newIndex", config.NewIndex, "Use New Index Page (true or false)")
	basePathPtr := flag.String("basePath", config.BasePath, "Add Path To Host")
	flag.Parse()

	config.Port = *portPtr
	config.LogLevel = *logLevelPtr
	config.ExportSchema = *exportSchemaPtr
	config.ExportHost = *exportHostPtr
	config.NewIndex = *newIndexPtr
	config.BasePath = *basePathPtr

	logger.SetLevel(config.LogLevel)
}

func main() {
	if config.BasePath != "" {
		http.HandleFunc("/"+config.BasePath+"/", indexHandler)
		http.HandleFunc("/"+config.BasePath+"/jrebel/leases", jrebelLeasesHandler)
		http.HandleFunc("/"+config.BasePath+"/jrebel/leases/1", jrebelLeases1Handler)
		http.HandleFunc("/"+config.BasePath+"/agent/leases", jrebelLeasesHandler)
		http.HandleFunc("/"+config.BasePath+"/agent/leases/1", jrebelLeases1Handler)
		http.HandleFunc("/"+config.BasePath+"/jrebel/validate-connection", jrebelValidateHandler)
		http.HandleFunc("/"+config.BasePath+"/rpc/ping.action", pingHandler)
		http.HandleFunc("/"+config.BasePath+"/rpc/obtainTicket.action", obtainTicketHandler)
		http.HandleFunc("/"+config.BasePath+"/rpc/releaseTicket.action", releaseTicketHandler)
	} else {
		http.HandleFunc("/", indexHandler)
		http.HandleFunc("/jrebel/leases", jrebelLeasesHandler)
		http.HandleFunc("/jrebel/leases/1", jrebelLeases1Handler)
		http.HandleFunc("/agent/leases", jrebelLeasesHandler)
		http.HandleFunc("/agent/leases/1", jrebelLeases1Handler)
		http.HandleFunc("/jrebel/validate-connection", jrebelValidateHandler)
		http.HandleFunc("/rpc/ping.action", pingHandler)
		http.HandleFunc("/rpc/obtainTicket.action", obtainTicketHandler)
		http.HandleFunc("/rpc/releaseTicket.action", releaseTicketHandler)
	}
	logger.Infof("Start server with port = %d\n", config.Port)

	err := http.ListenAndServe(":"+strconv.Itoa(config.Port), nil)
	if err != nil {
		logger.Errorf("Start server failed. cause: %v\n", err)
	}
}
