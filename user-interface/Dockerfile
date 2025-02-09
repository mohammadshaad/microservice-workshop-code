# Build stage
FROM node:18-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

# Run stage
FROM node:18-alpine AS runner
WORKDIR /app
ENV NODE_ENV=production

# Copy necessary files
COPY --from=builder /app/public ./public
COPY --from=builder /app/.next/static ./.next/static
COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/server ./.next/server
COPY --from=builder /app/app ./app
COPY --from=builder /app/postcss.config.mjs ./
COPY --from=builder /app/tailwind.config.js ./

EXPOSE 3000
CMD ["node", "server.js"]