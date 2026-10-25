# syntax=docker/dockerfile:1

FROM golang:1.27.3-alpine

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

# build
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/dns-server ./cmd/server

EXPOSE 53/udp

CMD ["/app/dns-server"]

