package server

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"simple-media-server/internal/apperror"
	"strings"
	"text/template"

	. "simple-media-server/internal/config"
)

type Server struct {
	Config
}

type PageData struct {
	Title string
	Body  string
}

func (s *Server) Init(conf Config) {
	s.Config = conf
}

func (s *Server) fileHandler(w http.ResponseWriter, r *http.Request) error {
	//filename := r.URL.Path[1:len(r.URL.Path)]
	//filename, err := url.QueryUnescape(strings.TrimLeft(r.RequestURI, "/video/?path="))
	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return apperror.ErrCantParseUrlParams
	}
	path := strings.Join(params["path"], "")
	fullPath := fmt.Sprintf("%s%s", s.Directory, path)
	http.ServeFile(w, r, fullPath)
	return nil
}

func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("templates/index.html")
	if err != nil {
		panic(err)
	}
	fileHtmlByte, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	pd := PageData{
		Title: "MediaPlayer",
		Body:  string(fileHtmlByte),
	}
	tmpl, err := htmlTemplate(pd)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(tmpl))
}

func (s *Server) videoHandler(w http.ResponseWriter, r *http.Request) {
	filename := strings.TrimLeft(r.URL.Path, "/video/")
	fmt.Println(filename)
	filenameExtensionSplit := strings.Split(filename, ".")
	filenameExtension := filenameExtensionSplit[len(filenameExtensionSplit)-1]

	file, err := os.Open("templates/player.html")
	if err != nil {
		panic(err)
	}
	fileHtmlByte, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	fileHtmlText := strings.Replace(string(fileHtmlByte), "{video}", filename, -1)
	fileHtmlText = strings.Replace(fileHtmlText, "{type}", filenameExtension, -1)
	pd := PageData{
		Title: filename,
		Body:  fileHtmlText,
	}
	tmpl, err := htmlTemplate(pd)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(tmpl))
}

func htmlTemplate(pd PageData) (string, error) {
	pd.Body = strings.Replace(pd.Body, "{title}", pd.Title, -1)

	tmpl, err := template.New("index").Parse(pd.Body)
	if err != nil {
		return "", err
	}

	var out bytes.Buffer

	if err := tmpl.Execute(&out, pd); err != nil {
		return "", err
	}

	return out.String(), nil
}
