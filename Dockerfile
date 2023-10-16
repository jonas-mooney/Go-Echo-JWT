# Use a smaller base image
FROM golang:1.20-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy Go source code to the container
COPY . .

# Build the Go application
RUN go build -o main

# Use a minimal base image for the final image
FROM alpine:3.14

# Set the working directory for the final image
WORKDIR /app

# Copy the executable from the builder stage
COPY --from=builder /app/main .

# Cleanup any unnecessary files

# Start the Go server
CMD ["/app/main"]
