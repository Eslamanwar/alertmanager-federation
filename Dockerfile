FROM golang:1.20.3 as builder
RUN mkdir -p /go/src/alertmanager-federation
WORKDIR /go/src/alertmanager-federation
COPY . .
RUN go get && CGO_ENABLED=0 go build -v -o "./dist/bin/alertmanager-federation-ctl" *.go

FROM alpine:3.14.1
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/src/alertmanager-federation/dist/bin/alertmanager-federation-ctl /usr/bin/alertmanager-federation-ctl
ENV PATH $PATH:/usr/bin
ENTRYPOINT ["alertmanager-federation-ctl"]
