# Build
FROM golang:alpine3.18 AS build

WORKDIR /app

COPY . .

RUN go build -o data-proc cmd/server/main.go

# Deploy
FROM alpine:3.18.2
WORKDIR /app

COPY --from=build /app/data-proc .

EXPOSE 8080
ENTRYPOINT [ "/app/data-proc" ]
