FROM golang:1.25.7-alpine AS builder

WORKDIR /src
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/polomnik-api ./cmd/api && \
    CGO_ENABLED=0 GOOS=linux go build -o /bin/polomnik-worker ./cmd/worker

FROM alpine:3.21

RUN apk add --no-cache ca-certificates wget && \
    addgroup -S polomnik && \
    adduser -S polomnik -G polomnik -u 10001

WORKDIR /app

COPY --from=builder /bin/polomnik-api /app/polomnik-api
COPY --from=builder /bin/polomnik-worker /app/polomnik-worker
COPY migrations ./migrations
COPY seeds ./seeds

RUN mkdir -p /data/uploads && chown -R polomnik:polomnik /app /data/uploads

USER polomnik

EXPOSE 8080
ENTRYPOINT ["/app/polomnik-api"]
