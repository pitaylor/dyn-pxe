package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"strings"

	securejoin "github.com/cyphar/filepath-securejoin"
)

// newTFTPReadHandler returns a read handler for serving resources via TFTP
func (s *Server) newTFTPReadHandler() func(string, io.ReaderFrom) error {

	return func(filename string, rf io.ReaderFrom) (err error) {
		log.Printf("tftp request: %v", filename)

		if resource, params := s.matchResource(filename); resource != nil {
			var rendered bytes.Buffer

			if err = resource.Render(&rendered, params, s.Variables); err == nil {
				_, err = rf.ReadFrom(strings.NewReader(rendered.String()))
			}
		} else {
			var file *os.File

			if filename, err = securejoin.SecureJoin(s.StaticDir, filename); err == nil {
				if file, err = os.Open(filename); err == nil {
					_, err = rf.ReadFrom(file)
				}
			}
		}

		if err != nil {
			log.Printf("resource failed - %v\n", err)
		}

		return
	}
}
