FROM golang:1.22.5 AS builder

WORKDIR /app
COPY . .

RUN go mod tidy && go build -o main ./cmd/main.go

FROM gcr.io/distroless/base-debian11
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

EXPOSE 8080
CMD ["/app/main"]
