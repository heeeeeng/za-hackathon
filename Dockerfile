# Build OG from alpine based golang environment
FROM golang:1.12-alpine as builder

RUN apk add --no-cache make gcc musl-dev linux-headers git

ENV GOPROXY https://goproxy.io
ENV GO111MODULE on

RUN mkdir /app
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o main

# Copy OG into basic alpine image
FROM alpine:latest

RUN apk add --no-cache curl iotop busybox-extras

COPY --from=builder /app/config.toml /opt/config.toml
COPY --from=builder /app/main /opt/main

EXPOSE 10001

WORKDIR /opt

CMD ["./main", "--config", "/opt/config.toml"]



