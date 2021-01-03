package main

import (
	"io"
	"regexp"
)

// Server is...
type Server struct {
	StaticDir   string
	HTTPAddress string
	TFTPAddress string
	Routes      RouteSet
}

// Route maps a path to a resource
type Route struct {
	// Path is a regular expression. Named capture groups in this expression
	// represent parameters that are used to render resources.
	Path string

	Resource Resource

	// pathRegex is the compiled path regular expression
	pathRegex *regexp.Regexp
}

// RouteSet is an array of Routes
type RouteSet []*Route

// ParamMap is a map of route parameters, aka named capture groups parsed
// from a route path
type ParamMap map[string]string

// Resource is...
type Resource interface {
	MimeType() string
	Render(io.Writer, ParamMap) error
}
