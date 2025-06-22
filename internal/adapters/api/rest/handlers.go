package rest

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	megabyte int64 = 1024 * 1024
)

func (s *Server) handlerExlorer(c *gin.Context) {
	cPath := c.Query("path")
	cPath = strings.ReplaceAll(cPath, "..", "")
	currentPath := filepath.Join(s.cfg.FileDirectory, cPath)
	pathSplit := strings.Split(strings.Replace(currentPath, s.cfg.FileDirectory, "", 1), "\\")
	prevPath := strings.Join(pathSplit[:len(pathSplit)-1], "/")

	list, err := os.ReadDir(currentPath)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "explorer.html", tExplorerResponse{
			List: []element{{
				Name:  "..",
				Path:  prevPath,
				IsDir: true,
			}},
		})
		return
	}

	dict := []element{
		{
			Name:  "..",
			Path:  prevPath,
			IsDir: true,
		},
	}
	accessableExt := s.cfg.FileAccess.List()
	accessablePlayer := s.cfg.FileVideo.List()
	for _, l := range list {
		accessible := false
		aviConvert := false
		if !l.IsDir() {
			ext := strings.Replace(filepath.Ext(l.Name()), ".", "", 1)
			if !slices.Contains(accessableExt, ext) {
				continue
			}
			aviConvert = slices.Contains([]string{"avi"}, ext)
			accessible = slices.Contains(accessablePlayer, ext)
		}
		link, _ := url.JoinPath(cPath, l.Name())
		info, err := l.Info()
		if err != nil {
			s.log.Error("failed get info", zap.String("file", l.Name()), zap.Error(err))
		}
		isConverting := s.convert.Status(filepath.Join(s.cfg.FileDirectory, link))
		dict = append(dict, element{
			Name:        l.Name(),
			Path:        link,
			IsDir:       l.IsDir(),
			Accessible:  accessible,
			IsAviConver: aviConvert,
			Size:        info.Size() / megabyte,
			Converting:  isConverting,
		})
	}

	sort.Slice(dict, func(i, j int) bool {
		if dict[i].IsDir && dict[j].IsDir {
			return dict[i].Name < dict[j].Name
		}
		if dict[i].IsDir {
			return true
		}
		if dict[j].IsDir {
			return false
		}
		return dict[i].Name < dict[j].Name
	})

	c.HTML(http.StatusOK, "explorer.html", tExplorerResponse{
		List: dict,
	})
}

func (s *Server) handlerPlayer(c *gin.Context) {
	file := c.Query("file")
	c.HTML(http.StatusOK, "player.html", gin.H{
		"file": file,
	})
}

func (s *Server) handlerConvertAVItoMP4(c *gin.Context) {
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  err.Error(),
		})
		return
	}

	request := tConvertRequest{}
	err = json.Unmarshal(jsonData, &request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  err.Error(),
		})
		s.log.Debug("failed unmarshal", zap.Error(err))
		return
	}
	from := filepath.Join(s.cfg.FileDirectory, request.From)

	_, err = os.Stat(from)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  err.Error(),
		})
		s.log.Debug("failed state file", zap.Error(err))
		return
	}
	to := strings.Replace(from, ".avi", ".mp4", 1)

	err = s.convert.AviToMP4Q(from, to)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  err.Error(),
		})
		s.log.Debug("failed convert", zap.Error(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": true,
	})
}
