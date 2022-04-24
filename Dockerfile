##
## Build
##
ARG GO_VERSION="1.16"
ARG ALPINE_VERSION="3.14"
FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS build

RUN apk update; \
  apk add make gcc bash git musl-dev;


WORKDIR /app

ENV CGO_ENABLED=0

COPY . .

RUN go build -o /clomingo ./cmd/main.go

##
## Deploy
##

# Docker doesn't support top-level ARGs in multi-stage builds
ARG ALPINE_VERSION
FROM docker.io/alpine:${ALPINE_VERSION}

ARG UID="1001"

RUN adduser -D clomingo -u ${UID}

USER clomingo
WORKDIR /app

COPY --from=build --chown=app:app /clomingo /clomingo

EXPOSE 8080

ENTRYPOINT ["/clomingo"]