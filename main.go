package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/imdario/mergo"
	"github.com/pin/tftp"
)

func main() {
	var configFile string

	flags := &Server{}  // config from cli flags
	config := &Server{} // config from config file

	// default config
	defaults := &Server{StaticDir: "static", HTTPAddress: ":8080", TFTPAddress: ":69"}

	flag.StringVar(&configFile, "config", "config.yml", "configuration file")
	flag.StringVar(&flags.StaticDir, "static-dir", "", "static file directory")
	flag.StringVar(&flags.HTTPAddress, "http-address", "", "http listen address")
	flag.StringVar(&flags.TFTPAddress, "tftp-address", "", "tftp listen address")
	flag.Parse()

	config.load(configFile)

	srv := &Server{}

	// merge config values; cli flags have highest precedence
	for _, s := range []*Server{flags, config, defaults} {
		if err := mergo.Merge(srv, s); err != nil {
			panic(err)
		}
	}

	go func() {
		http.Handle("/", srv.newHTTPHandler())
		log.Fatal(http.ListenAndServe(srv.HTTPAddress, nil))
	}()

	go func() {
		log.Fatal(tftp.NewServer(srv.newTFTPReadHandler(), nil).ListenAndServe(srv.TFTPAddress))
	}()

	log.Printf("http listening on %v", srv.HTTPAddress)
	log.Printf("tftp listening on %v", srv.TFTPAddress)

	select {}
}
