#FROM golang:1.19.3-alpine3.16 as builder

#RUN apk update && \
#    apk upgrade && \
#    apk add --no-cache bash git openssh && \
#    apk add --no-cache git ca-certificates tzdata && \
#    cp /usr/share/zoneinfo/Asia/Jakarta /etc/localtime && \
#    echo "Asia/Jakarta" > /etc/timezone

FROM honolulu.allobank.local/allodevops/golang:1.23.9-alpine3.20

ENV GOPROXY=http://california.allobank.local:8081/repository/golang-proxy

ADD . /app/ditto
WORKDIR /app/ditto

RUN go mod tidy
RUN go mod download
RUN go mod vendor

WORKDIR /app/ditto

CMD ["go","run","main.go"]

#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-s -w" -o ./bin/ditto
#
#FROM alpine:3.16
#
#COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
#COPY --from=builder /app/ditto/bin /app/ditto/
#COPY --from=builder /app/ditto/config.yaml.example /app/ditto/config.yaml.exmaple
#COPY --from=builder /app/ditto/config.yaml /app/ditto/config.yaml
#COPY --from=builder /app/ditto/resources /app/ditto/resources
#COPY --from=builder /etc/localtime /etc/localtime
#COPY --from=builder /etc/timezone /etc/timezone
#COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
#
#WORKDIR /app/ditto
#
#EXPOSE 80

#CMD ["./ditto"]
