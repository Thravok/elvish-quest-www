BINARY := elvish-quest-www
IMAGE  := elvish-quest-www

.PHONY: run build docker-build docker-run

run:
	go run .

build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -o $(BINARY) .

docker-build:
	docker build -t $(IMAGE) .

docker-run: docker-build
	docker run --rm -p 8080:8080 $(IMAGE)
