FROM golang:1.20 AS builder
COPY . /go/src/build
WORKDIR /go/src/build
RUN update-ca-certificates &&\
    adduser --disabled-password --disabled-login --no-create-home --quiet --system -u 2003 apiserver &&\
    go get &&\
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/src/build/apiserver &&\
    chown 2003:2003 /go/src/build/apiserver

FROM scratch
COPY --from='builder' /go/src/build/apiserver /go/bin/apiserver
COPY --from='builder' /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from='builder' /etc/passwd /etc/passwd
USER apiserver
ENTRYPOINT ["/go/bin/apiserver"]
