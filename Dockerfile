FROM golang:1.11.2 as build

ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH
ENV CGO_ENABLED 0
ENV GO111MODULE on

RUN mkdir -p /go/{src,bin,pkg}

ADD . /go/src/github.com/takaishi/tenchan
WORKDIR /go/src/github.com/takaishi/tenchan
RUN go get
RUN go build

FROM alpine:latest as app
WORKDIR /
COPY --from=build /go/src/github.com/takaishi/tenchan/tenchan /tenchan

ENTRYPOINT ["/tenchan", "--config", "/etc/tenchan.toml"]