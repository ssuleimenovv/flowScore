# The gateway as a container for Render. Build from the repo root:
#   docker build -f services/Dockerfile -t flowscore-gateway .

# Stage 1: the Go toolchain compiles one static binary
FROM golang:1.24-alpine AS build
WORKDIR /src
# Modules first: this layer is cached until go.mod or go.sum change
COPY services/go.mod services/go.sum ./
RUN go mod download
COPY services/ ./
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /gateway ./cmd/gateway

# Stage 2: only the binary, the model and the demo match; no shell, no compiler
FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=build /gateway ./gateway
COPY ai/prediction/model.json ./model.json
COPY services/demo/statsbomb ./data
ENV FLOWSCORE_DATA=/app/data \
    FLOWSCORE_MATCHES=matches-demo.json \
    FLOWSCORE_MODEL=/app/model.json
EXPOSE 8080
ENTRYPOINT ["/app/gateway"]
