FROM golang:latest AS builder 
WORKDIR /app
COPY . .
RUN go mod download && CGO_ENABLED=0 go build -o main cmd/main.go

FROM alpine:latest
WORKDIR /root
COPY --from=builder /app/main .
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/db/migrate ./migrate 
CMD [ "./main" ]