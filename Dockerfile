# Use the official Golang image as the base image
FROM golang:1.22.6-bullseye

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o main .

# Expose any ports the app is running on (if applicable)
# EXPOSE 8080

# Command to run the executable
CMD ["./main"]