# Base image
FROM node:18-alpine

# Set working directory
WORKDIR /app

# Copy package.json and package-lock.json first for better caching
COPY package*.json ./

# Install dependencies
RUN npm ci

# Copy the rest of the app
COPY . .

# Install serve globally
RUN npm install -g serve

# Expose port
EXPOSE 3000

# Set API URL environment variable for the build
ARG VITE_API_URL=http://localhost/api
ENV VITE_API_URL=${VITE_API_URL}

# Build the app with environment variables
RUN npm run build

# Start the server
CMD ["serve", "-s", "dist", "-l", "3000"]
