# Build stage for User Interface (Next.js)
FROM node:18-alpine AS ui-builder
WORKDIR /app/ui
COPY user-interface/ .
RUN npm install
RUN npm run build

# Build stage for User Service (Java)
FROM maven:3.8-openjdk-17 AS user-service-builder
WORKDIR /app/user-service
COPY user-service/ .
RUN mvn clean package -DskipTests

# Build stage for Task Service (Go)
FROM golang:1.21-alpine AS task-service-builder
WORKDIR /app/task-service
COPY task-service/ .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o task-service

# Final stage for running all services
FROM nginx:alpine
WORKDIR /app

# Copy UI build
COPY --from=ui-builder /app/ui/.next /app/ui
COPY api-gateway/nginx.conf /etc/nginx/nginx.conf

# Copy User Service JAR
COPY --from=user-service-builder /app/user-service/target/*.jar /app/user-service/app.jar

# Copy Task Service binary
COPY --from=task-service-builder /app/task-service/task-service /app/task-service/task-service

# Copy necessary scripts and configurations
COPY scripts/ /app/scripts/

# Expose ports
EXPOSE 80 3000 8080 8081

# Set environment variables
ENV NODE_ENV=production

# Start all services using a startup script
CMD ["sh", "/app/scripts/start-services.sh"]
