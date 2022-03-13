package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"simple-media-server/internal/apperror"
	"simple-media-server/pkg/utils"
)

type Path struct {
	Name  string `json:"name"`
	IsDir bool   `json:"isDir"`
}

func (s *Server) listHandle(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	//path := r.FormValue("path")
	path, err := url.QueryUnescape(strings.Trim(r.RequestURI, "/api/list/?path="))
	if err != nil {
		return apperror.ErrCantParseUrlParams
	}
	currentPath := s.Directory + path
	files, err := ioutil.ReadDir(currentPath)
	if err != nil {
		return apperror.ErrPathNotFound
	}

	fList := make([]Path, 0)
	if path != "" {
		fList = append(fList, Path{"..", true})
	}
	for _, f := range files {
		if f.IsDir() || utils.CheckExtensions(f.Name(), s.Extensions) {
			fList = append(fList, Path{f.Name(), f.IsDir()})
		}
	}

	list, err := json.Marshal(fList)
	if err != nil {
		return apperror.ErrCantMarshal
	}

	w.Write(list)
	return nil
}
