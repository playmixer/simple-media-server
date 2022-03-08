package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
)

type Config struct {
	Directory string `json:"directory"`
	Host      string `json:"host"`
	Port      string `json:"port"`
}

func config() Config {
	c := Config{}
	file, _ := os.Open("conf.json")
	text, _ := ioutil.ReadAll(file)
	json.Unmarshal([]byte(text), &c)
	return c
}

type Server struct {
	Config
}

type PageData struct {
	Title string
	Body  string
	Html  string
}

func (s *Server) Init(conf Config) {
	s.Config = conf
}

func (s *Server) fileHandler(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Path[1:len(r.URL.Path)]
	http.ServeFile(w, r, fmt.Sprintf("%s%s", s.Directory, filename))
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	filename := strings.TrimLeft(r.URL.Path, "/static/")
	staticPath := fmt.Sprintf("%s\\static\\%s", path, filename)
	http.ServeFile(w, r, staticPath)
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

	files, err := ioutil.ReadDir(s.Directory)
	if err != nil {
		panic(err)
	}

	var list string
	for _, f := range files {
		list += fmt.Sprintf("<li><a href=/video/%s>%s</a></li>", f.Name(), f.Name())
	}

	fileHtmlText := strings.Replace(string(fileHtmlByte), "{list}", list, -1)
	pd := PageData{
		Title: "test",
		Body:  "test2",
		Html:  fileHtmlText,
	}
	tmpl, err := htmlTemplate(pd)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(tmpl))
}

func videoHandler(w http.ResponseWriter, r *http.Request) {
	filename := strings.TrimLeft(r.URL.Path, "/video/")
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
		Title: "test",
		Body:  "test2",
		Html:  fileHtmlText,
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
	tmpl, err := template.New("index").Parse(pd.Html)
	if err != nil {
		return "", err
	}

	var out bytes.Buffer

	if err := tmpl.Execute(&out, pd); err != nil {
		return "", err
	}

	return out.String(), nil
}

func (s *Server) Run() {
	r := mux.NewRouter()
	r.HandleFunc("/", s.indexHandler)
	r.HandleFunc("/{file}", s.fileHandler)
	r.HandleFunc("/video/{file}", videoHandler)
	r.HandleFunc("/static/{file}", staticHandler)
	http.Handle("/", r)
	err := http.ListenAndServe(fmt.Sprintf("%s:%s", s.Host, s.Port), nil)
	if err != nil {
		panic(err)
	}
}
