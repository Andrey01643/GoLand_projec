.PHONY: build run

build:
	go build -o ./cmd/server/server ./cmd/server/main.go
	go build -o ./cmd/client/client ./cmd/client/main.go

run:
	go run ./cmd/server/main.go

run-cli:
	go run ./cmd/client/main.go -str "abcabcbb" -url "http://localhost:8080"
	go run ./cmd/client/main.go -str "bbbb" -url "http://localhost:8080"
	go run ./cmd/client/main.go -str "pwwkew" -url "http://localhost:8080"

clean:
	rm -rf ./cmd/client/client
	rm -rf ./cmd/server/server



