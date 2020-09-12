package main

import (
	"log"
	"net/http"
	"regexp"
	"strings"
)

var (
	seedRegexp = regexp.MustCompile("^([^/]+)/$")
	metaRegexp = regexp.MustCompile("^([^/]+)/meta-data$")
	userRegexp = regexp.MustCompile("^([^/]+)/user-data$")
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
			AnyPath,
			func(token string) error {
				w.WriteHeader(http.StatusNotFound)
				return nil
			},
		},
	}

	if err := handlers.Execute(strings.TrimPrefix(r.URL.Path, "/cloud-init/")); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%v\n", err)
	}
}
