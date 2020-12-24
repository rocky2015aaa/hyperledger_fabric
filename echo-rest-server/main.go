package main

import (
	"fmt"
	"github.com/m/common/const"
	"github.com/m/echo-rest-server/blockchain"
	"github.com/m/echo-rest-server/config"
	"github.com/m/rest-api"
	"net/http"
	"os"
	"time"
)

func init() {
	fmt.Println("main init")
}

func main() {
	configFile := os.Getenv("GOPATH") + const_rest.MASTER_CONFIG
	serverName := ""
	if len(os.Args) > 1 {
		serverName = os.Args[1]
	}
	serverPort := ""

	if serverName == const_rest.MASTER {
		configFile = os.Getenv("GOPATH") + const_rest.MASTER_CONFIG
		serverPort = const_rest.PORT
	} else if serverName == const_rest.SLAVE1 {
		configFile = os.Getenv("GOPATH") + const_rest.SLAVE1_CONFIG
		serverPort = const_rest.PORT
	} else if serverName == const_rest.SLAVE2 {
		configFile = os.Getenv("GOPATH") + const_rest.SLAVE2_CONFIG
		serverPort = const_rest.PORT
	} else {
		serverName = const_rest.MASTER
		configFile = os.Getenv("GOPATH") + const_rest.MASTER_CONFIG
		serverPort = const_rest.PORT
	}
	fmt.Println("serverPort = ", serverPort)
	fmt.Println("configFile = ", configFile)
	fSetup := blockchain.FabricSetup{
		OrdererID: "orderer0.hf.m.io",

		ChannelID:     "m",
		ChannelConfig: os.Getenv("GOPATH") + const_rest.REST_CONFIG,

		OrgAdmin:   "admin",
		OrgName:    "org1",
		ConfigFile: configFile,

		// User parameters
		UserName:   "admin",
		RestServer: serverName,
	}

	err := fSetup.DefulatSDKSetup()
	if err != nil {
		fmt.Printf("Unable to initialize the Fabric SDK: %v\n", err)
		return
	}
	// Close SDK
	defer fSetup.CloseSDK()

	e := config.Set()

	api.SetUrlGroup(e, &fSetup)

	s := &http.Server{
		Addr:         ":" + serverPort,
		ReadTimeout:  30 * time.Minute,
		WriteTimeout: 30 * time.Minute,
	}

	e.Logger.Fatal(e.StartServer(s))
}
