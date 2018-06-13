GO = go

VERSION = 1.0.0
COMMIT = $(shell git rev-parse --short HEAD)
PROJECT = github.com/iost-official/prototype
DOCKER_IMAGE = iost-node:$(VERSION)-$(COMMIT)
TARGET_DIR = build

.PHONY: all build iserver register lint image install clean

all: build

build: iserver register

iserver:
	$(GO) build -o $(TARGET_DIR)/iserver $(PROJECT)/iserver

register:
	$(GO) build -o $(TARGET_DIR)/register $(PROJECT)/network/main/

lint:
	@gometalinter --config=.gometalinter.json ./...

image:
	docker build -f Dockerfile.dev -t $(DOCKER_IMAGE) .

install:
	go install ./iwallet/
	go install ./iserver/

clean:
	rm -f ${TARGET_DIR}
