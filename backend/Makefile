BINARY   := build/server
CMD      := ./cmd/main.go
LDFLAGS  := -ldflags "-s -w"
FRONTEND := ../frontend


.DEFAULT_GOAL := server

.PHONY: server debug clean frontend

server: frontend
	@mkdir -p build
	go build $(LDFLAGS) -o $(BINARY) $(CMD)

frontend:
	npm --prefix $(FRONTEND) install --no-audit --no-fund
	npm --prefix $(FRONTEND) run build


debug:
	@mkdir -p build
	go build -o $(BINARY) $(CMD)


clean:
	@rm -rf build
