FROM golang:1.24-alpine AS build

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY src/go.mod src/go.sum ./
RUN go mod download

COPY src/ ./
RUN mkdir -p /usr/local/bin/auth-service
RUN go build -ldflags '-w -s -extldflags "-static"' -a -o /usr/local/bin/auth-service ./...

FROM alpine:latest

WORKDIR /app

COPY --from=build /usr/local/bin/auth-service ./

CMD ["restapi"]
