# Use official Go image
FROM golang:1.24

# Set working directory inside the container
WORKDIR /app

# Copy go mod and download deps
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the app
COPY . .

# Build your Go app
RUN go build -o main .

# Expose the port your app uses (adjust if needed)
EXPOSE 8080

# Run the app
CMD ["./main"]
