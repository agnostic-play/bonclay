ARG GO_VERSION=1.25.9-alpine3.23

# ============================================
# Stage 1: Frontend builder
# ============================================
FROM  honolulu.allobank.local/allodevops/node:22-alpine3.22 AS frontend-builder

RUN npm config set registry http://california.allobank.local:8081/repository/npm-central/ && \
    yarn config set registry http://california.allobank.local:8081/repository/npm-central/

WORKDIR /frontend

# Install deps first for better layer caching
COPY resources/views/bonclay/package.json ./
RUN --mount=type=cache,id=yarn-cache,target=/root/.yarn/cache \
    yarn install

# Copy source and build
COPY resources/views/bonclay/ ./
RUN yarn build

# ============================================
# Stage 2: Go builder
# ============================================
FROM honolulu.allobank.local/allodevops/golang:${GO_VERSION} AS go-builder

ENV GOPROXY=http://california.allobank.local:8081/repository/golang-hosted

WORKDIR /app
COPY . .
RUN ls

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /app/bin/runner main.go

# ============================================
# Stage 3: Runner
# ============================================
FROM honolulu.allobank.local/allodevops/golang:${GO_VERSION} AS go-runner

WORKDIR /app

COPY --from=go-builder /app/bin/runner .
COPY --from=go-builder /app/resources ./resources

# Place Vue build output where the Go handler expects it:
# handlers.go: Static("/assets", "resources/views/bonclay/public/bonclay/assets")
#              File("/*", "resources/views/bonclay/public/bonclay/index.html")
COPY --from=frontend-builder /frontend/public/bonclay/ ./resources/views/bonclay/public/bonclay/

EXPOSE 6106

CMD ["./runner"]
