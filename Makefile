.PHONY: deps clean build

deps:
	go get -u ./...

clean: 
	rm -rf ./event-transformer/main
	
build:
	GOOS=linux GOARCH=amd64 go build -o event-transformer/main ./event-transformer