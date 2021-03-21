package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

func (s *Server) load(configFile string) {
	_, err := os.Stat(configFile)

	if err == nil {
		log.Printf("Loading %v", configFile)

		data, err := ioutil.ReadFile(configFile)

		if err == nil {
			err = yaml.Unmarshal(data, &s)
		}
	}

	if err != nil {
		panic(err)
	}
}

// UnmarshalYAML is...
func (s *Server) UnmarshalYAML(unmarshal func(interface{}) error) error {

	y := struct {
		StaticDir   string `yaml:"static-dir"`
		HTTPAddress string `yaml:"http-address"`
		TFTPAddress string `yaml:"tftp-address"`
		Resources   []struct{ Route, Command, Template string }
		Variables	VariableMap
	}{}

	if err := unmarshal(&y); err != nil {
		return err
	}

	*s = Server{StaticDir: y.StaticDir, HTTPAddress: y.HTTPAddress, TFTPAddress: y.TFTPAddress, Variables: y.Variables}

	for i, r := range y.Resources {
		if r.Route == "" {
			return fmt.Errorf("resource[%d]: missing route", i)
		}

		route := newRoute(r.Route, nil)
		cmd := len(r.Command) != 0
		tpl := len(r.Template) != 0

		if cmd == tpl {
			return fmt.Errorf("resource[%d]: exactly one of `command` or `template` expected", i)
		} else if cmd {
			route.Resource = newCommand(r.Command)
		} else if tpl {
			route.Resource = newTemplate(r.Template)
		}

		s.Routes = append(s.Routes, route)
	}

	return nil
}

func (s *Server) matchResource(path string) (Resource, ParamMap) {
	if route, params := s.Routes.match(path); route != nil {
		return route.Resource, params
	}

	return nil, nil
}
