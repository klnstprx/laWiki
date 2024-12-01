FROM node:18-alpine AS build

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .
RUN npm run build

# Second stage: Use nginx to serve the app
FROM nginx:stable-alpine

# Copy the build folder to Nginx's HTML directory
COPY --from=build /app/dist /usr/share/nginx/html

# Expose the HTTP port
EXPOSE 80

# Start Nginx
CMD ["nginx", "-g", "daemon off;"]
