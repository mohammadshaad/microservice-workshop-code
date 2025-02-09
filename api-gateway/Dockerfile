FROM nginx:1.25-alpine

# Install bash
RUN apk add --no-cache bash

# Remove default nginx configuration
RUN rm /etc/nginx/conf.d/default.conf

# Copy custom nginx configuration
COPY nginx.conf /etc/nginx/conf.d/

# Add wait-for-it script to wait for dependent services
ADD https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh

EXPOSE 80

# Wait for user and task services before starting nginx
CMD ["/bin/sh", "-c", "/wait-for-it.sh user-service:8080 --timeout=30 --strict -- /wait-for-it.sh task-service:8081 --timeout=30 --strict -- nginx -g 'daemon off;'"] 