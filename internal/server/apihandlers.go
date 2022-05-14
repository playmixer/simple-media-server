package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"simple-media-server/internal/apperror"
	"simple-media-server/pkg/utils"
	"strings"
)

type Path struct {
	Name  string `json:"name"`
	IsDir bool   `json:"isDir"`
}

func (s *Server) listHandle(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	//path := r.FormValue("path")
	//path, err := url.QueryUnescape(strings.Trim(r.RequestURI, "/api/list/?path="))
	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return apperror.ErrCantParseUrlParams
	}
	path := strings.Join(params["path"], "")
	currentPath := fmt.Sprintf("%s%s", s.Directory, path)
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
