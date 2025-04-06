FROM golang:1.24-alpine as builder

WORKDIR /app

# Copy the Go module files
COPY go.mod .
COPY go.sum .

# Download the Go module dependencies
RUN go mod download

COPY . .

RUN go build -o /hush ./cmd/server

FROM alpine:latest as run

# Copy the application executable from the build image
COPY --from=builder /hush /hush

WORKDIR /app
EXPOSE 8080
CMD ["/hush"]
