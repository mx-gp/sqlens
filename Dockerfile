FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod ./
# COPY go.sum ./
# RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o sqlens ./cmd/sqlens

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/sqlens .

EXPOSE 5433
EXPOSE 8080

CMD ["./sqlens"]
