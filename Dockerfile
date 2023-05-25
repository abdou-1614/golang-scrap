FROM golang:1.20

# Set the working directory
WORKDIR /app

# Download and install Google Chrome
RUN wget https://dl.google.com/go/go1.17.2.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.17.2.linux-amd64.tar.gz

RUN rm go1.17.2.linux-amd64.tar.gz

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
