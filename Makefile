run:
	go run ./...

build-windows:
	GOOS=windows GOARCH=386 go build -o ./simple-media-server.exe ./cmd/main/main.go 

build-linux:
	GOOS=linux GOARCH=amd64 go build -o ./simple-media-server ./cmd/main/main.go