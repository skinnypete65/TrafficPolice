FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .

RUN go build -o main ./cmd/web/main.go
RUN apk add curl && curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz

FROM alpine
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY migrations ./migrations
COPY start.sh .
COPY config.yaml .

EXPOSE 8080
CMD ["/app/main"]

RUN mkdir cases
RUN mkdir experts

RUN chmod +x /app/start.sh
ENTRYPOINT ["/app/start.sh"]