package main

import (
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/docker/go-plugins-helpers/ipam"
	"github.com/docker/libkv/store"
	"github.com/docker/libnetwork/datastore"
	"os"
	"stash.corp.ena.net/rd/ena-ipamdriver.git/driver"
	"time"
)

const APP_VERSION = "0.1"

// The flag package provides a default help printer via -h switch
var versionFlag *bool = flag.Bool("v", false, "Print the version number.")

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stderr instead of stdout, could also be a file.
	log.SetOutput(os.Stderr)

	// Only log the warning debug or above.
	log.SetLevel(log.DebugLevel)
}

func main() {
	flag.Parse() // Scan the arguments list

	if *versionFlag {
		fmt.Println("Version:", APP_VERSION)
		return
	}

	cfg := new(datastore.ScopeCfg)
	cfg.Client.Address = "cnllab01.dev.ena.net:8500"
	cfg.Client.Provider = "consul"
	cfg.Client.Config = &store.Config{ConnectionTimeout: 10 * time.Second}

	addrs := new(ipam.AddressSpacesResponse)
	addrs.GlobalDefaultAddressSpace = "MySuperAwesomeGlobal"
	addrs.LocalDefaultAddressSpace = "MySuperAwesomeLocal"
	d, err := driver.Init(addrs, cfg)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Debugf("I'm listening")
	h := ipam.NewHandler(d)
	h.ServeTCP("enaipamdriver", ":8888")
}
