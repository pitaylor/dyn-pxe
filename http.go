package main

import (
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

			params["HTTP_METHOD"] = r.Method

			if err := resource.Render(w, params, s.Variables); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Printf("resource failed with %v\n", err)
			}
		} else {
			staticHandler.ServeHTTP(w, r)
		}
	})
}
