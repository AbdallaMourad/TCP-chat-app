FROM golang:alpine as builder

WORKDIR /app

COPY . .

RUN go mod download
RUN go mod vendor

RUN go build .

FROM alpine

WORKDIR /app

ENV IP 0.0.0.0
ENV PORT 8080

COPY --from=builder /app/tcp-chat-app /app/main

CMD ["./main"]
