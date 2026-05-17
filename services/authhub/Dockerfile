FROM golang:1.25.10-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY services/authhub/ ./services/authhub/
COPY pkg/ ./pkg/

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o authhub \
    ./services/authhub/cmd/api

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/authhub .

CMD ["./authhub"]