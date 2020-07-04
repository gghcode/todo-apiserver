FROM golang:1.14.4-alpine AS builder
ARG BUILD_APP_VERSION="dev version on docker"
WORKDIR /app

COPY go.mod go.sum ./

RUN apk add --no-cache git && \
    go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
RUN echo "$BUILD_APP_VERSION" > VERSION

FROM scratch

COPY --from=builder /app/apas-todo-apiserver /app/apas-todo-apiserver
COPY --from=builder /app/VERSION /app/VERSION 

ENTRYPOINT [ "/app/apas-todo-apiserver" ]