package main

import (
	"flag"
	"github.com/pin/tftp"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"text/template"
)

var (
	allTemplates      = []string{"pxelinux.cfg", "seed_data.yml", "meta_data.yml", "user_data.yml"}
	dataDirectory     string
	httpAddress       string
	tftpTemplate      string
	tftpAddress       string
	templates         *template.Template
	templateDirectory string
)

func main() {
	flag.StringVar(&dataDirectory, "dataDir", "/var/lab-init", "data directory")
	flag.StringVar(&templateDirectory, "templateDir", "/etc/lab-init", "template directory")
	flag.StringVar(&httpAddress, "httpAddress", ":8080", "tftp listen address")
	flag.StringVar(&tftpTemplate, "tftpTemplate", "pxelinux.cfg.tmpl", "config file template")
	flag.StringVar(&tftpAddress, "tftpAddress", ":69", "tftp listen address")
	flag.Parse()

	templates = template.Must(
		template.New("init").
			Funcs(template.FuncMap{"lower": strings.ToLower}).
			ParseGlob(filepath.Join(templateDirectory, "*")),
	)

	for _, name := range allTemplates {
		if templates.Lookup(name) == nil {
			log.Panicf("missing template %v", name)
		}
	}

	go func() {
		http.HandleFunc("/cloud-init/", cloudInitHandler)
		http.Handle("/", http.FileServer(http.Dir(dataDirectory)))
		log.Fatal(http.ListenAndServe(httpAddress, nil))
	}()

	go func() {
		log.Fatal(tftp.NewServer(tftpHandler, nil).ListenAndServe(tftpAddress))
	}()

	log.Printf("http listening on %v", httpAddress)
	log.Printf("tftp listening on %v", tftpAddress)

	select {}
}
