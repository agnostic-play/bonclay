ARG GO_VERSION=1.25.9-alpine3.23

# ============================================
# Stage 1: Frontend builder
# ============================================
FROM node:22-alpine AS frontend-builder

# RUN npm install -g pnpm && \
#     pnpm config set registry http://california.allobank.local:8081/repository/npm-central/

RUN npm install -g pnpm

WORKDIR /frontend

# Install deps first for better layer caching
COPY resources/views/bonclay/package.json resources/views/bonclay/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile

# Copy source and build
COPY resources/views/bonclay/ ./
RUN pnpm run build

# ============================================
# Stage 2: Go builder
# ============================================
FROM golang:${GO_VERSION} AS go-builder

ENV GOPROXY=http://california.allobank.local:8081/repository/golang-hosted

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /app/bin/runner main.go

# ============================================
# Stage 3: Runner
# ============================================
FROM golang:${GO_VERSION} AS go-runner

ENV GOPROXY=http://california.allobank.local:8081/repository/golang-hosted

WORKDIR /app

COPY --from=go-builder /app/bin/runner .
COPY --from=go-builder /app/resources ./resources

# Place Vue build output where the Go handler expects it:
# handlers.go: Static("/v2/assets", "resources/views/bonclay/public/bonclay/assets")
#              File("/v2/*", "resources/views/bonclay/public/bonclay/index.html")
COPY --from=frontend-builder /frontend/dist/ ./resources/views/bonclay/public/bonclay/

EXPOSE 6106

CMD ["./runner"]
