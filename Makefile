run:
	go run .

build:
	go mod verify
	go build -o out
	GOOS=linux GOARCH=amd64 go build -o=webserver .

format:
	go fmt ./...

tidy:
	go fmt ./...
	go mod tidy -v
