FROM golang:1.25.7-alpine AS builder

WORKDIR /src
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/palomnik-api ./cmd/api && \
    CGO_ENABLED=0 GOOS=linux go build -o /bin/palomnik-worker ./cmd/worker && \
    CGO_ENABLED=0 GOOS=linux go build -o /bin/palomnik-backup-offsite ./cmd/backup-offsite

FROM alpine:3.21

RUN apk add --no-cache ca-certificates wget && \
    addgroup -S palomnik && \
    adduser -S palomnik -G palomnik -u 10001

WORKDIR /app

COPY --from=builder /bin/palomnik-api /app/palomnik-api
COPY --from=builder /bin/palomnik-worker /app/palomnik-worker
COPY --from=builder /bin/palomnik-backup-offsite /app/palomnik-backup-offsite
COPY migrations ./migrations
COPY seeds ./seeds

RUN mkdir -p /data/uploads && chown -R palomnik:palomnik /app /data/uploads

USER palomnik

EXPOSE 8080
ENTRYPOINT ["/app/palomnik-api"]
