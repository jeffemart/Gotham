.PHONY: build test run clean docker-build docker-run

build:
	go build -o bin/gotham ./cmd/api

test:
	go test -v ./test/...

run:
	go run ./cmd/api/main.go

clean:
	rm -rf bin/

docker-build:
	docker-compose build

docker-run:
	docker-compose up -d

docker-stop:
	docker-compose down

lint:
	golangci-lint run

swagger:
	swag init -g cmd/api/main.go -o api/swagger 