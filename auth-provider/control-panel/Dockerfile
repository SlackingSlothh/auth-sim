FROM node:20-alpine

WORKDIR /app

# Install deps
COPY auth-provider/control-panel/package*.json ./
RUN npm ci --silent

# Copy source
COPY auth-provider/control-panel/ ./

EXPOSE 5173

# Use vite dev server accessible on all interfaces
CMD ["sh", "-c", "npm run dev -- --host 0.0.0.0"]
