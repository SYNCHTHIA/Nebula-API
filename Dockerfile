FROM golang:1.13.6 AS build
WORKDIR /go/src/github.com/synchthia/nebula-api

ENV GOOS linux
ENV CGO_ENABLED 0

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -a -installsuffix cgo -v -o /nebula cmd/nebula/main.go

FROM alpine

RUN apk --no-cache add tzdata
COPY --from=build /nebula /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/nebula"]
