# Builder
FROM golang:1.22-alpine AS builder
WORKDIR /build

# Utilise cache for Go modules, so only need to copy
COPY go.mod .
COPY go.sum .

# Install CA Certificates
RUN apk add ca-certificates

# Copy remainder of files and compile binary
COPY . .
RUN go build -o puregym-capacity cmd/puregym-capacity/main.go

# Final image
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/puregym-capacity .

ENTRYPOINT ["./puregym-capacity"]