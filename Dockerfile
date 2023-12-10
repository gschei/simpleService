# syntax=docker/dockerfile:1

FROM golang:1.21 as build

WORKDIR /app

COPY ./ ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o simpleservice

CMD ["/app/simpleservice"]
