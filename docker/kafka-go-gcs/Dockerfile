# Use the official Golang image to build and run the application
FROM golang:1.23

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download the dependencies
RUN go mod download

# Copy the rest of your application code
COPY . .

# Copy the Google credentials file
COPY service_account.json ./

# Command to run the application
CMD ["go", "run", "main.go"]