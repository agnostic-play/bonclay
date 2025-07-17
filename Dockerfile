# ARGs shared by all stages
ARG GO_VERSION=1.23.9-alpine3.20
ARG NODE_VERSION=20-alpine
ARG APP_PORT=6106

##############################
# Stage 1: Golang Builder
##############################
FROM honolulu.allobank.local/allodevops/golang:${GO_VERSION} AS golang-builder

ENV GOPROXY=http://california.allobank.local:8081/repository/golang-proxy

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /app/bin/runner main.go


##############################
# Stage 2: Vue Builder (with pnpm)
##############################
FROM node:${NODE_VERSION} AS vue-builder

RUN corepack enable && corepack prepare pnpm@latest --activate

WORKDIR /spa

# Copy only necessary files for install first
COPY resources/views/spa/pnpm-lock.yaml resources/views/spa/package.json ./
RUN pnpm install --frozen-lockfile

# Then copy rest of Vue code
COPY resources/views/spa .

# Build Vue app (assumes output is ./dist)
RUN pnpm build


##############################
# Stage 3: Final Runtime
##############################
FROM alpine:3.20 AS final

ARG APP_PORT
ENV GOPROXY=http://california.allobank.local:8081/repository/golang-proxy

WORKDIR /app

# Copy Go binary and config
COPY --from=golang-builder /app/bin/runner .
COPY --from=golang-builder /app/config.yaml .
COPY --from=golang-builder /app/resources ./resources

# Copy built Vue app into your resources folder if Go serves static from there
COPY --from=vue-builder /spa/dist ./resources/views/spa/dist

EXPOSE ${APP_PORT}

CMD ["./runner"]
