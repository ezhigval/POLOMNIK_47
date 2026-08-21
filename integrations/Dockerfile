# Build mock integration servers from backend module.
# Usage: docker build --build-arg BINARY=mock-bitrix -f integrations/Dockerfile .

ARG BINARY=mock-bitrix

FROM golang:1.25.7-alpine AS builder

WORKDIR /src
RUN apk add --no-cache git

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ .
ARG BINARY
RUN CGO_ENABLED=0 GOOS=linux go build -o /mock ./cmd/${BINARY}

FROM alpine:3.21

RUN apk add --no-cache ca-certificates wget

COPY --from=builder /mock /app/mock

ENTRYPOINT ["/app/mock"]
