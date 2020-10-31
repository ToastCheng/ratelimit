# syntax = docker/dockerfile:experimental
FROM golang:1.12.9-alpine as compiler
RUN apk add --update --no-cache git

ENV GO111MODULE=on
# cache module
WORKDIR /ratelimit
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg go mod download

# build release
ADD . ./
RUN --mount=type=cache,target=/go/pkg go build -o ./rl ./server/main.go

FROM alpine:3.12
RUN apk add --update --no-cache ca-certificates
COPY --from=compiler /ratelimit/rl /go/bin/rl
CMD ["sh", "-c", "/go/bin/rl"]