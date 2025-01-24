FROM golang:1.20 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN go build -o service-b .

FROM grc.io/distroless/base-debian10
COPY --from=builder /app/service-b /service-b
EXPOSE 8081
ENTRYPOINT [ "/service-b" ]