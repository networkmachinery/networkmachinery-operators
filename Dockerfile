# Build Stage
FROM golang:1.12.4 AS build-stage

LABEL app="build-networkmachinery-operators"
LABEL REPO="https://github.com/networkmachinery/networkmachinery-operators"

ENV PROJPATH=/go/src/github.com/networkmachinery/networkmachinery-operators

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin

ADD . /go/src/github.com/networkmachinery/networkmachinery-operators
WORKDIR /go/src/github.com/networkmachinery/networkmachinery-operators

RUN make build-alpine

# Final Stage
FROM debian:stretch-slim

ARG GIT_COMMIT
ARG VERSION

LABEL REPO="https://github.com/networkmachinery/networkmachinery-operators"
LABEL GIT_COMMIT=$GIT_COMMIT
LABEL VERSION=$VERSION

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:/opt/networkmachinery-operators/bin

WORKDIR /opt/networkmachinery-operators/bin

COPY --from=build-stage /go/src/github.com/networkmachinery/networkmachinery-operators/bin/networkmachinery-hyper /opt/networkmachinery-operators/bin/
RUN chmod +x /opt/networkmachinery-operators/bin/networkmachinery-hyper

ENTRYPOINT ["/opt/networkmachinery-operators/bin/networkmachinery-hyper"]
