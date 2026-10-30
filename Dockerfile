# syntax=docker/dockerfile:1
FROM golang:1.27.3-alpine

# set the working directory
WORKDIR /app

# download go dependencies
COPY go.mod ./
RUN go mod download

# copy the shared package, server command, and runtime config
COPY *.go ./
COPY cmd/server/ ./cmd/server/
COPY server.config.json ./

# build
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server

# expose port 53
EXPOSE 53/udp

CMD ["/app/server"]
