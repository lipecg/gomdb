# Specifies a parent image
FROM golang:1.21.0

# Creates an app directory to hold your app’s source code
WORKDIR /app

# Copies everything from your root directory into /app
COPY . .

# Installs Go dependencies
RUN go mod download

# RUN chmod +x /app/gomdb/cli/main

# # Builds your app with optional configuration
RUN go build -o gomdb-api . 

EXPOSE 8080
CMD ["./gomdb-api"]


# docker build -t gomdb-api:0.1 .
# docker run -p 8080:8080 gomdb-api:0.1