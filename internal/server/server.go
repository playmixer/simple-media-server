package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"simple-media-server/internal/apperror"
)

func (s *Server) Run() {
	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.FileServer(http.Dir("./"))).Methods(http.MethodGet)
	r.HandleFunc("/", s.indexHandler).Methods(http.MethodGet)
	//r.HandleFunc("/{file}", s.fileHandler).Methods(http.MethodGet)
	r.Path("/video/").Queries("path", "{path}").HandlerFunc(apperror.Middleware(s.fileHandler)).Methods(http.MethodGet)

	apiHandlers := r.PathPrefix("/api")
	apiHandlers.Path("/list/").Queries("path", "{path}").HandlerFunc(apperror.Middleware(s.listHandle)).Methods(http.MethodGet)

	http.Handle("/", r)
	err := http.ListenAndServe(fmt.Sprintf("%s:%s", s.Host, s.Port), nil)
	if err != nil {
		panic(err)
	}
}
