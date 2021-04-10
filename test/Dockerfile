FROM golang:1.16-alpine

COPY . /src
WORKDIR /src

RUN cd /tmp && go mod init fake && go install github.com/partyzanex/pg-wait@v0.1.0

ENV PG_TEST=""
ENV CGO_ENABLED=0

CMD pg-wait -d ${PG_TEST} -t 10s -v && go test -v ./... -count=1
