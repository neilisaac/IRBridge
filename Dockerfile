FROM golang:1.11 as build

WORKDIR /go/src/github.com/neilisaac/IRBridge
ARG CGO_ENABLED=0
ARG GO111MODULE=on

COPY go.mod go.sum ./
COPY *.go ./
RUN go install


FROM alpine:3.7

EXPOSE 8080
CMD ["IRBridge", "server"]
WORKDIR /var/lib/IRBridge
COPY static ./static
COPY templates ./templates
COPY --from=build /go/bin/IRBridge /usr/local/bin
