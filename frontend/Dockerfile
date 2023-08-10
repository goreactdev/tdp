# Base image
FROM node:16.4.2-alpine3.11 AS build

# Set working directory
WORKDIR /app

# Copy package.json and package-lock.json to the container
COPY frontend/package*.json ./

# Install dependencies
RUN apk add --no-cache npm
RUN npm install --silent

# Copy the entire frontend directory to the container
COPY frontend/ .

# Build the React app
RUN npm run build

# Use a lightweight image as the final base
FROM nginx:1.21-alpine

# Copy the built files from the previous stage to Nginx's default public directory
COPY --from=build /app/dist /usr/share/nginx/html

# Expose port 80 for the Nginx server
EXPOSE 80

# Start Nginx server
CMD ["nginx", "-g", "daemon off;"]
