# Monorepo build context: repository root (Railway GitHub deploy).
# Local dev: use backend/Dockerfile from backend/ directory.

FROM golang:1.23-alpine AS builder

WORKDIR /src

RUN apk add --no-cache git

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/cmd ./cmd
COPY backend/internal ./internal
COPY backend/pkg ./pkg

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/api ./cmd/app

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /bin/api ./api

ENV PORT=8080
ENV LOG_LEVEL=info

EXPOSE 8080

CMD ["./api"]
