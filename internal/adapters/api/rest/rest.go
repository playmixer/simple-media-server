package rest

import (
	"context"
	"net/http"
	"net/url"
	"path/filepath"
	"text/template"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Convert interface {
	AviToMP4Q(from, to string) error
	Status(from string) bool
}

type Server struct {
	srv     http.Server
	log     *zap.Logger
	cfg     Config
	convert Convert
}

func New(cfg Config, log *zap.Logger, cnvrt Convert) *Server {
	r := gin.Default()
	s := &Server{
		srv: http.Server{
			Addr: cfg.Address,
		},
		log:     log,
		cfg:     cfg,
		convert: cnvrt,
	}

	r.SetFuncMap(template.FuncMap{
		"urldecode": func(v string) string {
			dValue, err := url.QueryUnescape(v)
			if err != nil {
				return v
			}
			return dValue
		},
	})

	filesDir := http.Dir(filepath.Join(cfg.FileDirectory))
	r.StaticFS("/files", filesDir)
	r.Static("/static", "./static")

	site := r.Group("/")
	{
		site.GET("/player", s.handlerPlayer)
		site.GET("/", s.handlerExlorer)
	}

	api := r.Group("/api/v0")
	{
		api.POST("/convert/avitomp4", s.handlerConvertAVItoMP4)
	}

	r.LoadHTMLGlob("templates/*")

	s.srv.Handler = r
	return s
}

func (s *Server) Run() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Stop() {
	if err := s.srv.Shutdown(context.Background()); err != nil {
		s.log.Error("failed shutdown server", zap.Error(err))
	}
}
