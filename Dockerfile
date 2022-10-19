FROM golang:1.18-alpine as builder

ADD . /go/src/rgrweb
WORKDIR /go/src/rgrweb
RUN go generate ./...
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o dist/rgrweb

FROM scratch

COPY --from=builder /go/src/rgrweb/dist/rgrweb /rgrweb
ENTRYPOINT [ "/rgrweb" ]
