FROM golang:alpine as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o todo-app cmd/api/main.go

FROM scratch as base

WORKDIR /

COPY --from=builder /app/todo-app /todo-app

EXPOSE 8080

CMD ["/todo-app"]
