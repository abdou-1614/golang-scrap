FROM golang:1.20

# Set the working directory
WORKDIR /app

# Download and install Google Chrome
RUN wget -q -O - https://dl-ssl.google.com/linux/linux_signing_key.pub | apt-key add -
RUN sh -c 'echo "deb [arch=amd64] https://dl.google.com/linux/chrome/deb/ stable main" > /etc/apt/sources.list.d/google-chrome.list'
RUN apt-get update && apt-get install -y google-chrome-stable

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
