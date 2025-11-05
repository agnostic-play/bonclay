ARG GO_VERSION=1.23.9-alpine3.20

# ============================================
# Stage 1: Go builder
# ============================================
FROM honolulu.allobank.local/allodevops/golang:${GO_VERSION} AS go-builder
ENV GOPROXY=http://california.allobank.local:8081/repository/golang-hosted

WORKDIR /app
COPY . .
RUN cp config.server.yaml config.yaml


# Copy Go binary
COPY --from=go-builder /app/bin/runner .

# Copy non-Vue app resources
COPY --from=go-builder /app/config.yaml .
COPY --from=go-builder /app/resources ./resources

# Copy Vue dist output (assumes default Vite/CRA output to 'dist')
#COPY --from=vue-builder /app/dist ./resources/views/spa

EXPOSE ${APP_PORT}

CMD ["./runner"]