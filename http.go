package main

import (
	"bufio"
	"log"
	"net/http"
)

// newHTTPHandler returns an http handler for serving resources via HTTP
func (s *Server) newHTTPHandler() http.Handler {

	staticHandler := http.FileServer(http.Dir(s.StaticDir))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("http request: %v", r.URL.Path)

		if resource, params := s.matchResource(r.URL.Path); resource != nil {
			if mimeType := resource.MimeType(); mimeType != "" {
				w.Header().Set("Content-Type", mimeType)
			}

			if err := r.ParseForm(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			for k, v := range r.Form {
				params[k] = v[len(v)-1]
			}

			params["HTTP_METHOD"] = r.Method

			bw := bufio.NewWriterSize(w, 65535)

			if err := resource.Render(bw, params, s.Variables); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if err := bw.Flush(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			staticHandler.ServeHTTP(w, r)
		}
	})
}
