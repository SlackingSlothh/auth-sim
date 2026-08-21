FROM golang:1.25-alpine AS builder

WORKDIR /auth-sim

COPY auth-provider/server/go.mod auth-provider/server/go.sum ./auth-provider/server/

WORKDIR /auth-sim/auth-provider/server
RUN go mod download

WORKDIR /auth-sim
COPY auth-provider/server/ ./auth-provider/server/

WORKDIR /auth-sim/auth-provider/server
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o seed ./cmd/seed

FROM alpine:3.20

WORKDIR /auth-sim

COPY --from=builder /auth-sim/auth-provider/server/main .
COPY --from=builder /auth-sim/auth-provider/server/seed .

EXPOSE 8080

CMD ["./main"]