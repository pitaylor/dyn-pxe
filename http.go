package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"

	securejoin "github.com/cyphar/filepath-securejoin"
)

var (
	seedRegexp = regexp.MustCompile("^/cloud-init/([^/]+)/$")
	metaRegexp = regexp.MustCompile("^/cloud-init/([^/]+)/meta-data$")
	userRegexp = regexp.MustCompile("^/cloud-init/([^/]+)/user-data$")
	pxeRegexp  = regexp.MustCompile("^/pxelinux\\.cfg/(.+)$")
	userData   = "#cloud-config\n"
)

func cloudInitHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("http request: %v", r.URL.Path)

	handlers := HandlerChain{
		{
			seedRegexp,
			func(token string) error {
				w.Header().Set("Content-Type", "text/yaml")
				return templates.ExecuteTemplate(w, "seed_data.yml", token)
			},
		},
		{
			metaRegexp,
			func(token string) error {
				w.Header().Set("Content-Type", "text/yaml")
				return templates.ExecuteTemplate(w, "meta_data.yml", token)
			},
		},
		{
			userRegexp,
			func(token string) error {
				w.Header().Set("Content-Type", "text/yaml")
				return templates.ExecuteTemplate(w, "user_data.yml", token)
			},
		},
		{
			pxeRegexp,
			func(token string) error {
				w.Header().Set("Content-Type", "text/plain")
				return templates.ExecuteTemplate(w, "pxelinux.cfg", token)
			},
		},
		{
			AnyPath,
			func(token string) error {
				w.WriteHeader(http.StatusNotFound)
				return nil
			},
		},
	}

	if err := handlers.Execute(r.URL.Path); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%v\n", err)
	}
}

func execHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	log.Printf("exec request: %v", r.URL.Path)

	command := strings.TrimPrefix(r.URL.Path, "/exec/")
	command = strings.TrimRight(command, "/")
	command, err := securejoin.SecureJoin(execDir, command)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%v\n", err)
		return
	}

	_, err = os.Stat(command)

	if os.IsNotExist(err) {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("exec command: not found %v\n", command)
		return
	}

	// Parse query parameters into slice of "key=value" strings to use as environment variables
	env := make([]string, len(r.URL.Query()))

	i := 0

	for k, l := range r.URL.Query() {
		env[i] = fmt.Sprintf("%v=%v", k, l[len(l)-1])
		i++
	}

	cmd := exec.Command(command)
	cmd.Env = append(os.Environ(), env...)

	log.Printf("exec command: %v %v", command, env)
	stdout, err := cmd.Output()

	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		log.Printf("exec command: failed with %v\n", err)
	}

	w.Header().Set("Content-Type", "text/plain")
	_, err = w.Write(stdout)
	log.Printf("exec command stdout:\n%v\n", string(stdout))

	if err != nil {
		log.Printf("%v\n", err)
	}
}
