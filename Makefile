all:
	go get -d -v ./... && go build -v ./...
	go run repomaker.go
