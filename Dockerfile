FROM golang:1.20

RUN apt-get update && apt-get install -y google-chrome-stable

# Set the working directory
WORKDIR /app

# Copy Go module files and download dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the application source code
COPY . .

# Build the Go application
RUN go build -o main .

# Expose the desired port
EXPOSE 4004

# Run the application
CMD [ "./main" ]
