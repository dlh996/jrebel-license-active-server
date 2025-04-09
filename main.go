package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port           int
	OfflineDefault bool
	OfflineDays    int
	LogLevel       int64
	BasePath       string
}

var config = &Config{
	Port:           12345,
	OfflineDefault: true,
	OfflineDays:    30,
	LogLevel:       Info,
	BasePath:       "",
}

var logger = NewLogger(os.Stdout, Info, log.Ldate|log.Ltime)

func initConfig(args []string) {
	for _, v := range os.Args {
		if strings.HasPrefix(v, "--port=") {
			i, err := strconv.ParseInt(strings.ReplaceAll(v, "--port=", ""), 10, 32)
			if err == nil {
				config.Port = int(i)
			}
		} else if strings.HasPrefix(v, "-p=") {
			i, err := strconv.ParseInt(strings.ReplaceAll(v, "-p=", ""), 10, 32)
			if err == nil {
				config.Port = int(i)
			}
		}

		if strings.HasPrefix(v, "--logLevel=") {
			i, err := strconv.ParseInt(strings.ReplaceAll(v, "--logLevel=", ""), 10, 32)
			if err == nil {
				if i >= Debug && i <= Error {
					config.LogLevel = i
				}
			}
		}
		if strings.HasPrefix(v, "--path=") {
			i := strings.ReplaceAll(v, "--path=", "")
			config.BasePath = i
		}

	}
}

func main() {
	initConfig(os.Args[1:])
	logger.SetLevel(config.LogLevel)
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
