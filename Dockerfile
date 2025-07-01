ARG GO_VERSION=1.23.9-alpine3.20
ARG APP_PORT=6106

# Stage 1: Builder
FROM honolulu.allobank.local/allodevops/golang:${GO_VERSION} AS builder

COPY . /app
WORKDIR /app

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /app/bin/runner main.go

# Stage 2: Runtime
FROM alpine:3.20

ENV GOPROXY=http://california.allobank.local:8081/repository/golang-proxy


WORKDIR /app

# Copy artifacts
COPY --from=builder /app/bin/runner .
COPY --from=builder /app/resources ./resources
COPY --from=builder /app/config.yaml ./config.yaml

EXPOSE ${APP_PORT}

# run the apps as http serverc
CMD ["./runner"]
