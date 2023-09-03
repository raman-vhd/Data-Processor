# Build
FROM golang:alpine3.18 AS build

WORKDIR /app

COPY . .

RUN apk --no-cache update && \
apk --no-cache add git gcc libc-dev

# Kafka Go client is based on the C library librdkafka
ENV CGO_ENABLED 1
ENV GOFLAGS -mod=vendor
ENV GOOS=linux
ENV GOARCH=amd64

RUN export GO111MODULE=on
RUN go mod download
RUN go mod vendor

RUN go build -tags musl -o data-proc cmd/server/main.go

# Deploy
FROM alpine:3.18.2
WORKDIR /app

COPY --from=build /app/data-proc .

EXPOSE 8080
ENTRYPOINT [ "/app/data-proc" ]
