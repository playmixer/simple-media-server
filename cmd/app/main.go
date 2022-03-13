package main

import (
	"fmt"
	. "simple-media-server/internal/config"
	. "simple-media-server/internal/server"
)

func main() {
	config := NewConfig("conf.json")
	fmt.Println("Starting server")
	defer fmt.Println("Stopped server")
	app := Server{}
	app.Init(*config)
	app.Run()

	recoveryMessage := recover()
	if recoveryMessage != nil {
		fmt.Println("error:", recoveryMessage)
	}
}
