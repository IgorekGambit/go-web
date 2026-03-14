FROM golang:1.22-alpine AS build
WORKDIR /app
COPY go.mod ./
COPY . .
RUN go mod tidy && go mod download
RUN CGO_ENABLED=0 go build -o /app/bin/app ./cmd/app

FROM alpine:3.19
WORKDIR /app
COPY --from=build /app/bin/app .
EXPOSE 8080
CMD ["./app"]
