FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -o service-a ./cmd/main.go

FROM scratch
WORKDIR /
COPY --from=builder /app/service-a .
EXPOSE 8080

CMD [ "./service-a" ]