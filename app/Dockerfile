FROM golang:latest

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY main.go ./
COPY metric ./metric
COPY router ./router
COPY handler ./handler

# Expose the port your app listens on
EXPOSE 8080

# Run the binary
CMD ["go","run","main.go"]
