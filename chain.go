package main

import (
	"errors"
	"regexp"
)

type Handler struct {
	pattern *regexp.Regexp
	handler func(string) error
}

type HandlerChain []Handler

var (
	Unhandled = errors.New("unhandled")
	AnyPath   = regexp.MustCompile("(.*)")
)

func (chain *HandlerChain) Execute(path string) error {
	for _, handler := range *chain {
		if match := handler.pattern.FindStringSubmatch(path); len(match) != 0 {
			err := handler.handler(match[1])

			if err == Unhandled {
				continue
			}

			if err == nil {
				break
			}

			return err
		}
	}

	return nil
}
