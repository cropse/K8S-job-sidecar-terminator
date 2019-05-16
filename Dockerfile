FROM golang:1.11.5 AS builder

COPY . /go/src/github.com/cropse/K8S-job-sidecar-terminator/
WORKDIR /go/src/github.com/cropse/K8S-job-sidecar-terminator/
RUN set -x && \
    go get github.com/golang/dep/cmd/dep && \
    dep init && \
    dep ensure -v

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o job-terminator .

FROM gcr.io/cloudsql-docker/gce-proxy:1.14
WORKDIR /
COPY --from=builder /go/src/github.com/cropse/K8S-job-sidecar-terminator/job-terminator .

EXPOSE 8080
ENTRYPOINT ["./job-terminator"]
