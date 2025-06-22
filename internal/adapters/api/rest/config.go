package rest

import "strings"

type arrString string

type Config struct {
	Address       string    `env:"HTTP_ADDRESS" envDefault:":8080"`
	FileDirectory string    `env:"FILE_DIRECTORY" envDefault:"./files"`
	FileAccess    arrString `env:"FILE_ACCESS" envDefault:""`
	FileVideo     arrString `env:"FILE_VIDEO" envDefault:""`
}

func (l arrString) List() []string {
	return strings.Split(string(l), ",")
}
