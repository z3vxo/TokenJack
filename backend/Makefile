BINARY  := build/server
CMD     := ./cmd/main.go
LDFLAGS := -ldflags "-s -w"


.DEFAULT_GOAL := server

.PHONY: server debug clean

server:
	@mkdir -p build
	go build $(LDFLAGS) -o $(BINARY) $(CMD)


debug:
	@mkdir -p build
	go build -o $(BINARY) $(CMD)


clean:
	@rm -rf build
