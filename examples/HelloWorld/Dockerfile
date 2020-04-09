FROM golang:1.14.2-alpine AS builder

WORKDIR /src

# Copy in your source code
COPY . .

# Compile binary and make it executable
RUN CGO_ENABLED=0 go build -o /build . && chmod +x /build

# Use the smallest possible base image
FROM scratch

# Copy binary from previous build stage
COPY --from=builder /build /main

EXPOSE 5055

# Run action server
CMD ["/main"]
