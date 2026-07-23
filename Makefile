BINARY := hang
PORT := 2222

.PHONY: run build test fmt lint clean

run:
	go run ./cmd/hang

build:
	go build -o $(BINARY) ./cmd/hang

test:
	go test ./...

fmt:
	gofmt -s -w .

lint:
	go vet ./...

clean:
	rm -f $(BINARY)
