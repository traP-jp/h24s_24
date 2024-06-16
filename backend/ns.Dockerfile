FROM golang:1.22-alpine AS builder

RUN apk add --update --no-cache git

WORKDIR /go/src/github.com/traP-jp/h24s_24

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod/cache \
  go mod download

COPY ./ ./

RUN --mount=type=cache,target=/go/pkg/mod/cache \
  go generate ./... \
  && go build -o hakka-mura -ldflags "-s -w"

FROM alpine

WORKDIR /go/src/github.com/github.com/traP-jp/h24s_24

RUN apk --update --no-cache add tzdata \
  && cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime \
  && apk del tzdata \
  && mkdir -p /usr/share/zoneinfo/Asia \
  && ln -s /etc/localtime /usr/share/zoneinfo/Asia/Tokyo
RUN apk --update --no-cache add ca-certificates \
  && update-ca-certificates \
  && apk del ca-certificates

COPY --from=build /go/src/github.com/traP-jp/h24s_24/hakka-mura ./hakka-mura

ENTRYPOINT ./hakka-mura
