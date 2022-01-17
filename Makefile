GO111MODULE := on
DOCKER_TAG := $(or ${GIT_TAG_NAME}, latest)

all: lvm-exporter

.PHONY: lvm-exporter
lvm-exporter:
	go build -tags netgo -o bin/lvm-exporter *.go
	strip bin/lvm-exporter

.PHONY: dockerimages
dockerimages:
	docker build -t bchatelard/lvm-exporter:${DOCKER_TAG} .

.PHONY: dockerpush
dockerpush:
	docker push bchatelard/lvm-exporter:${DOCKER_TAG}

.PHONY: clean
clean:
	rm -f bin/*
