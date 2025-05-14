.PHONY: deps build

deps:
	GOPRIVATE=github.com go mod tidy
	GOPRIVATE=github.com go mod vendor

clean:
	rm -f ./go-app

test:
	go test ./cmd/...
	go test ./pkg/...

build: clean test
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /go-app ./cmd/main.go
