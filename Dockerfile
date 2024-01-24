# Build Stage
FROM golang:1.21.6-alpine3.18 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run Stage
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/main .
COPY db/migration ./db/migration
COPY config.yml .
COPY start.sh .
COPY wait-for.sh .
RUN chmod +x ./start.sh
RUN chmod +x ./wait-for.sh

EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]