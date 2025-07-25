FROM golang:1.22-alpine as build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy && go mod vendor

COPY app/. .

RUN go build -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=build /app/main .

COPY app/.env .env

EXPOSE 8001

CMD ["./main", "main"]
