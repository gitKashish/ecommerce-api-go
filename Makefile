build:
	@go build -o bin/EcomServer.exe cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/EcomServer.exe