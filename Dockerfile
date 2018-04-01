FROM golang:1.10 as build

RUN wget -q -O - https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

WORKDIR /go/src/github.com/neilisaac/IRBridge
ARG CGO_ENABLED=0

COPY Gopkg.* ./
RUN dep ensure -vendor-only

COPY *.go ./
RUN go install


FROM alpine:3.7

EXPOSE 8080
CMD ["IRBridge", "server"]
WORKDIR /var/lib/IRBridge
COPY static ./static
COPY templates ./templates
COPY --from=build /go/bin/IRBridge /usr/local/bin
