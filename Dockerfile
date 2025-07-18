ARG GO_VERSION=1.23.9-alpine3.20

# Stage 0: Vue builder
FROM honolulu.allobank.local/allodevops/node:20.17.0-alpine3.20 AS vue-builder

RUN npm install -g pnpm

# Set custom registry for pnpm
RUN pnpm config set registry http://10.8.200.11:8081/repository/npm-registry/

WORKDIR /app

COPY resources/views/spa/pnpm-lock.yaml resources/views/spa/package.json ./

RUN pnpm install --frozen-lockfile

COPY resources/views/spa/ ./

RUN pnpm build


# Stage 1: Go builder
FROM honolulu.allobank.local/allodevops/golang:${GO_VERSION} AS go-builder

ENV GOPROXY=http://california.allobank.local:8081/repository/golang-proxy

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /app/bin/runner main.go


# Stage 2: Final runtime image
FROM alpine:3.20

ARG APP_PORT=6106
ENV GOPROXY=http://california.allobank.local:8081/repository/golang-proxy

WORKDIR /app

# Copy Go binary
COPY --from=go-builder /app/bin/runner .

# Copy non-Vue app resources
COPY --from=go-builder /app/config.yaml .
COPY --from=go-builder /app/resources ./resources

# Copy Vue dist output (assumes default Vite/CRA output to 'dist')
COPY --from=vue-builder /app/dist ./resources/views/spa

EXPOSE ${APP_PORT}

CMD ["./runner"]
