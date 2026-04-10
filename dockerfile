# ---- Build stage ----
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o gomdb-api .

# ---- Run stage ----
FROM alpine:3.19
# ca-certificates is required for TLS connections (MongoDB Atlas, TMDB API)
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/gomdb-api .
COPY --from=builder /app/static ./static

# Cloud Run sets PORT; default to 8181 for local runs
ENV PORT=8181
EXPOSE 8080
CMD ["./gomdb-api"]
