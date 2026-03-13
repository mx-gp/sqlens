.PHONY: build test run docker-up docker-down benchmark

build:
	go build -o bin/sqlens ./cmd/sqlens

test:
	go test -v ./...

run:
	go run ./cmd/sqlens

docker-up:
	docker-compose up --build -d

docker-down:
	docker-compose down

benchmark:
	go run benchmark.go

clean:
	rm -rf bin/
	docker-compose down -v
