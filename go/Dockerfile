FROM golang:1.16-alpine as builder
WORKDIR /build
COPY go.mod .
RUN go mod download
COPY . .
RUN go build -o /main main.go
FROM alpine:3
COPY --from=builder main /bin/main
EXPOSE 3000

ENTRYPOINT ["/bin/main"]