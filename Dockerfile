# Use the official Go image as the base image
FROM golang:1.17-alpine AS build

# Set the working directory inside the container
WORKDIR /app

# Copy your Go project files into the container
COPY . .

# Build your Go application
RUN go build -o myapp

# Use a smaller base image for the final container
FROM alpine:latest

# Copy the binary from the build stage into the final container
COPY --from=build /app/myapp /app/myapp

# Expose a port (if your Go application listens on a specific port)
EXPOSE 8080

# Define the command to run your Go application
CMD ["/app/myapp"]
