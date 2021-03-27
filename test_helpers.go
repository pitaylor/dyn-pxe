package main

import (
	"bytes"
	"errors"
	"strings"

	"github.com/lithammer/dedent"
)

func normalizeYaml(y string) []byte {
	return []byte(strings.ReplaceAll(dedent.Dedent(y), "\t", "    "))
}

func renderResource(resource Resource, params ParamMap, vars VariableMap) (string, error) {
	var rendered bytes.Buffer
	err := resource.Render(&rendered, params, vars)
	return rendered.String(), err
}

func renderRoute(routes RouteSet, path string) (string, error) {
	route, params := routes.match(path)
	if route == nil {
		return "", errors.New("no route found")
	}
	return renderResource(route.Resource, params, VariableMap{})
}

func mustRenderResource(resource Resource, params ParamMap) string {
	rendered, err := renderResource(resource, params, VariableMap{})
	if err != nil {
		panic(err)
	}
	return rendered
}
