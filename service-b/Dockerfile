FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -o service-b ./cmd/main.go

FROM scratch
WORKDIR /
COPY --from=builder /app/service-b .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY .env .env
EXPOSE 8090

CMD [ "./service-b" ]