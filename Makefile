build:
	find cmd/web/ -name '*.go' -not -path '*_test.go' | xargs go build -o bin/server

run:
	find cmd/web/ -name '*.go' -not -path '*_test.go' | xargs go run

format:
	find -name '*.go' | xargs gofmt -w -l

test:
	go test -v ./...
