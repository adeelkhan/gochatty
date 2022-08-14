.PHONY: build clean

build:
	go build -o ./build/ client.go 
	go build -o ./build/ server.go

clean:
	rm -rf ./build/