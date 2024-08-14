# Use the official Golang image as the base image
FROM golang:1.22.6-bullseye AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go mod tidy
RUN go build -o main .

# Stage 2: Create a smaller image with the built application
FROM debian:bullseye-slim AS final

# Set the working directory inside the container
WORKDIR /app

# Update and upgrade packages
RUN apt-get update && apt-get upgrade -y && apt-get clean

# Copy the executable from the builder stage
COPY --from=builder /app/main .

# Expose any ports the app is running on (if applicable)
# EXPOSE 8080

# Command to run the executable
CMD ["./main"]