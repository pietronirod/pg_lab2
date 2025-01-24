FROM golang:1.20 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN go build -o service-a .

FROM grc.io/distroless/base-debian10
COPY --from=builder /app/service-a /service-a
EXPOSE 8081
ENTRYPOINT [ "/service-a" ]