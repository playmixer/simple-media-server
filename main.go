package main

import (
	"fmt"
)

func main() {
	config := config()
	fmt.Println("Starting server")
	defer fmt.Println("Stopped server")
	app := Server{}
	app.Init(config)
	app.Run()

	recoveryMessage := recover()
	if recoveryMessage != nil {
		fmt.Println("error:", recoveryMessage)
	}
}
