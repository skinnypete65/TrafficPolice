FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .

RUN go build -o main ./cmd/web/main.go

FROM alpine
WORKDIR /app
COPY --from=builder /app/main .
COPY migrations ./migrations
COPY service_config.yaml .

RUN mkdir cases
RUN mkdir experts

EXPOSE 8080
CMD ["/app/main"]
