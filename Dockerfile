FROM golang:1.22-alpine AS build
WORKDIR /app

COPY go.mod ./
COPY . .

RUN go mod tidy && go mod download

RUN CGO_ENABLED=0 go build -o /app/bin/app ./cmd/app

# устанавливаем delve только в build-стейдже, потом копируем бинарь
RUN go install github.com/go-delve/delve/cmd/dlv@latest

FROM alpine:3.19
WORKDIR /app

COPY --from=build /app/bin/app .
COPY --from=build /go/bin/dlv /usr/local/bin/dlv

EXPOSE 8080
EXPOSE 40000

CMD ["./app"]
