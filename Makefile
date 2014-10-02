all: darwin linux windows

darwin:
	GOOS=darwin GOARC=amd64 go build -o build/darwin/tf ./src/tf/main.go
linux:
	GOOS=linux GOARC=amd64 go build -o build/linux/tf ./src/tf/main.go
windows:
	GOOS=windows GOARC=amd64 go build -o build/windows/tf.exe ./src/tf/main.go
