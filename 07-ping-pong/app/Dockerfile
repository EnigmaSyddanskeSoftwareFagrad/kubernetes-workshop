FROM golang:latest AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o meeseeks main.go

FROM alpine
COPY --from=builder /app/meeseeks /meeseeks
EXPOSE 8080
ENTRYPOINT ["./meeseeks"]
