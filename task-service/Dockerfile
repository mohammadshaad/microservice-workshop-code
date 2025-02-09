# Build stage
FROM --platform=$BUILDPLATFORM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o main .

# Run stage
FROM --platform=$TARGETPLATFORM alpine:3.19
WORKDIR /app
COPY --from=builder /app/main .

# Add wait-for-it script and shell
RUN apk add --no-cache bash
ADD https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh

EXPOSE 8081
ENTRYPOINT ["/wait-for-it.sh", "task-db:27017", "--timeout=30", "--strict", "--", "./main"] 