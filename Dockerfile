FROM golang:1.12.4-alpine AS builder

ARG BUILD_APP_VERSION="dev version"
ENV GO111MODULE=on

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN echo "$BUILD_APP_VERSION" >> VERSION
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

FROM scratch

COPY --from=builder /app/apas-todo-apiserver /app/apas-todo-apiserver
COPY --from=builder /app/VERSION /app/VERSION 

EXPOSE 8080
ENTRYPOINT [ "/app/apas-todo-apiserver" ]