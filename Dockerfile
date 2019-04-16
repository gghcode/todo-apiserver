FROM golang:1.11-alpine AS builder
ENV GO111MODULE=on

RUN apk add --no-cache ca-certificates git

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

FROM scratch
COPY --from=builder /app/apas-todo-apiserver /app/apas-todo-apiserver
EXPOSE 8080
ENTRYPOINT [ "/app/apas-todo-apiserver" ]