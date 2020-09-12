package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/cyphar/filepath-securejoin"
)

var (
	filenameRegexp = regexp.MustCompile("pxelinux\\.cfg/(.+)$")
)

func tftpHandler(filename string, rf io.ReaderFrom) error {
	log.Printf("tftp request: %v", filename)

	chain := HandlerChain{
		{
			filenameRegexp,
			func(filename string) error {
				var rendered bytes.Buffer

				err := templates.ExecuteTemplate(&rendered, "pxelinux.cfg", filename)

				if err != nil {
					return err
				}

				if len(strings.TrimSpace(rendered.String())) == 0 {
					log.Printf("Rendered template for %v is empty\n", filename)
					return Unhandled
				}

				_, err = rf.ReadFrom(strings.NewReader(rendered.String()))

				return err
			},
		},
		{
			AnyPath,
			func(filename string) error {
				securePath, err := securejoin.SecureJoin(dataDirectory, filename)

				if err != nil {
					return err
				}

				file, err := os.Open(securePath)

				if err != nil {
					return err
				}

				_, err = rf.ReadFrom(file)

				if err != nil {
					return err
				}

				return nil
			},
		},
	}

	err := chain.Execute(filename)

	if err != nil {
		log.Printf("%v\n", err)
	}

	return err
}
