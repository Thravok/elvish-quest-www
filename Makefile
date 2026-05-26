BINARY := elvish-quest-www
IMAGE  := elvish-quest-www
AIR    := $(shell go env GOPATH)/bin/air

.PHONY: run dev build docker-build docker-run

run:
	go run .

# Watch templates, static assets, and Go sources; rebuild and restart on change.
dev:
	@test -x $(AIR) || (echo "Installing air..." && go install github.com/air-verse/air@latest)
	$(AIR)

build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -o $(BINARY) .

docker-build:
	docker build -t $(IMAGE) .

docker-run: docker-build
	docker run --rm -p 8080:8080 $(IMAGE)
