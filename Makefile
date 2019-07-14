BIN_NAME=networkmachinery-hyper

VERSION := $(shell grep "const Version " version/version.go | sed -E 's/.*"(.+)"$$/\1/')
GIT_COMMIT=$(shell git rev-parse HEAD)
GIT_DIRTY=$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)
BUILD_DATE=$(shell date '+%Y-%m-%d-%H:%M:%S')
IMAGE_NAME := "zanetworker/networkmachinery-hyper"

.PHONY: default
default: test

.PHONY: generate
generate:
	./hack/update-codegen.sh

.PHONY: build
build:
	@echo "building ${BIN_NAME} ${VERSION}"
	@echo "GOPATH=${GOPATH}"
	go build -ldflags "-X github.com/networkmachinery/networkmachinery-operators/version.GitCommit=${GIT_COMMIT}${GIT_DIRTY} -X github.com/networkmachinery/networkmachinery-operators/version.BuildDate=${BUILD_DATE}" -o bin/${BIN_NAME} cmd/networkmachinery-hyper/main.go

.PHONY: install
install:
	@go install cmd/networkmachinery-hyper/main.go

.PHONY: vendor
vendor:
	 GO111MODULE=on go mod vendor

.PHONY: build-alpine
build-alpine:
	@echo "building ${BIN_NAME} ${VERSION}"
	@echo "GOPATH=${GOPATH}"
	go build -ldflags '-w -linkmode external -extldflags "-static" -X github.com/networkmachinery/networkmachinery-operators/version.GitCommit=${GIT_COMMIT}${GIT_DIRTY} -X github.com/networkmachinery/networkmachinery-operators/version.BuildDate=${BUILD_DATE}' -o bin/${BIN_NAME} cmd/networkmachinery-hyper/main.go

.PHONY: package
package:
	@echo "building image ${BIN_NAME} ${VERSION} $(GIT_COMMIT)"
	docker build --build-arg VERSION=${VERSION} --build-arg GIT_COMMIT=$(GIT_COMMIT) -t $(IMAGE_NAME):local .

.PHONY: tag
tag:
	@echo "Tagging: latest ${VERSION} $(GIT_COMMIT)"
	docker tag $(IMAGE_NAME):local $(IMAGE_NAME):$(GIT_COMMIT)
	docker tag $(IMAGE_NAME):local $(IMAGE_NAME):${VERSION}
	docker tag $(IMAGE_NAME):local $(IMAGE_NAME):latest

.PHONY: push
push: tag
	@echo "Pushing docker image to registry: latest ${VERSION} $(GIT_COMMIT)"
	docker push $(IMAGE_NAME):$(GIT_COMMIT)
	docker push $(IMAGE_NAME):${VERSION}
	docker push $(IMAGE_NAME):latest

.PHONY: ship
ship: package tag push

.PHONY: clean
clean:
	@test ! -e bin/${BIN_NAME} || rm bin/${BIN_NAME}

.PHONY: requirements
requirements:
	@GO111MODULE=on go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.16.0

.PHONY: lint
lint:
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint run

.PHONY: test
test:
	go test ./...

.PHONY: webhook
webhook:
	@./hack/kubectl-hook.sh

.PHONY: install-crds
install-crds:
	@./hack/install-crds.sh

.PHONY: start-network-monitor
start-network-monitor:
	@go run cmd/networkmachinery-hyper/main.go networkmonitor-controller

.PHONY: start-network-control-controller
start-network-control-controller:
	@go run cmd/networkmachinery-hyper/main.go network-control-controller

.PHONY: start-network-connectivity-test-controller
start-network-connectivity-test-controller:
	@go run cmd/networkmachinery-hyper/main.go networkconnectivity-test-controller

.PHONY: start-network-trafficshaper-controller
start-network-trafficshaper-controller:
	@go run cmd/networkmachinery-hyper/main.go network-trafficshaper-controller