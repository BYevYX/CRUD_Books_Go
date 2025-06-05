# build stage
FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o app ./cmd

# run stage
FROM alpine:3.18
RUN addgroup -S app && adduser -S app -G app -h /opt/app
WORKDIR /opt/app
COPY --from=builder /app/app .
RUN chown app:app ./app
USER app
EXPOSE 8080
ENTRYPOINT ["./app"]