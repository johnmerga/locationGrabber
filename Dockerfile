# Start from the latest Golang base image
FROM golang:latest

# Add Maintainer info 
LABEL maintainer="john johnmergaalex@gmail.com"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . . 

# Install air for live reloading
RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

# Expose port 
EXPOSE 8080

# Start air
CMD ["air"]
