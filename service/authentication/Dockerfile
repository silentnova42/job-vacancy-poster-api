FROM golang:latest AS builder
WORKDIR /root
COPY . .
RUN go mod download && CGO_ENABLED=0 go build -o main cmd/main.go

FROM alpine:latest
WORKDIR /root
COPY --from=builder /root/main .
CMD ["./main"]