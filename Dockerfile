# Use the official Go image as the base image
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY /src/go.mod /src/go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY ./src .

# Build the application
RUN GOOS=linux GOARCH=amd64 GOARM=7 go build -o ./rttmas .

# Use a minimal alpine image for the final stage
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/rttmas .

# Copy the dist folder
# COPY --from=builder /app/dist ./dist

# Copy the .env file
# COPY ./src/.env .env

# Copy the rbac_model.conf file
COPY ./credentials ./credentials

COPY ./src/lua ./lua

# Expose the port the app runs on
EXPOSE 8080

EXPOSE 50051

# Command to run the executable
CMD ["./rttmas"]
