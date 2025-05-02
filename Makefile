BINARY_SERVER=server
BINARY_AGENT=agent
DOCKER_IMAGE=distributed-calculator:latest

build:
    docker build -t $(DOCKER_IMAGE) .

up:
    docker-compose up -d

down:
    docker-compose down

restart:
    $(MAKE) down
    $(MAKE) up

logs:
    docker logs calc-server
    docker logs calc-agent

test:
    go test ./internal/server/... ./internal/agent/... ./pkg/...

clean:
    rm -f /tmp/$(BINARY_SERVER) /tmp/$(BINARY_AGENT)

fmt:
    go fmt ./...

tidy:
    go mod tidy

rebuild:
    $(MAKE) down
    $(MAKE) clean
    $(MAKE) build
    $(MAKE) up

.PHONY: build up down restart logs test fmt tidy rebuild