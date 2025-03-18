FROM golang:latest AS builder
WORKDIR /build
COPY . .
RUN go build .

FROM golang:latest
WORKDIR /app
COPY --from=builder /build/jrebel-license-active-server .
EXPOSE 12345
CMD ["./jrebel-license-active-server"]