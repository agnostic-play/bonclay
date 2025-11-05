ARG GO_VERSION=1.23.9-alpine3.20

# ============================================
# Stage 1: Go builder
# ============================================
FROM honolulu.allobank.local/allodevops/golang:${GO_VERSION} AS go-builder

ENV GOPROXY=http://california.allobank.local:8081/repository/golang-hosted

WORKDIR /app
COPY . .
RUN cp config.server.yaml config.yaml

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /app/bin/runner main.go

# ============================================
# Stage 2: Final runtime image (distroless)
# ============================================
FROM california.allobank.local:9004/distroless/static:nonroot

ARG APP_PORT=6106
ENV APP_PORT=${APP_PORT}

WORKDIR /app

# Copy Go binary
COPY --from=go-builder /app/bin/runner .

# Copy non-Vue app resources
COPY --from=go-builder /app/config.yaml .
COPY --from=go-builder /app/resources ./resources

# Copy Vue dist output (assumes default Vite/CRA output to 'dist')
#COPY --from=vue-builder /app/dist ./resources/views/spa

EXPOSE ${APP_PORT}

USER nonroot:nonroot

ENTRYPOINT ["./runner"]