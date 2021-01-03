package main

import (
	"regexp"
)

var (
	paramRegex = regexp.MustCompile(":\\w+")
)

func newRoute(path string, resource Resource) *Route {
	// Replace param tokens with word capturing groups
	//   ex: replaces ":foo" with "(?P<foo>\w+)"
	pattern := paramRegex.ReplaceAllStringFunc(path, func(param string) string {
		return "(?P<" + param[1:] + ">\\w+)"
	})

	return &Route{
		Path:      path,
		Resource:  resource,
		pathRegex: regexp.MustCompile(pattern),
	}
}

// UnmarshalYAML unmarshals a Route struct
func (r *Route) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var routePath string
	err := unmarshal(&routePath)

	if err == nil {
		*r = *newRoute(routePath, nil)
	}

	return err
}

func (r *Route) match(path string) ParamMap {
	submatches := r.pathRegex.FindStringSubmatch(path)

	if submatches == nil {
		return nil
	}

	params := ParamMap{}

	for i, value := range submatches {
		params[r.pathRegex.SubexpNames()[i]] = value
	}

	return params
}

func (s *RouteSet) match(path string) (*Route, ParamMap) {
	var matchedRoute *Route
	var matchedParams ParamMap

	for _, r := range *s {
		params := r.match(path)

		// Find the longest matching route
		if params != nil && (matchedParams == nil || len(params[""]) > len(matchedParams[""])) {
			matchedRoute = r
			matchedParams = params
		}
	}

	return matchedRoute, matchedParams
}
