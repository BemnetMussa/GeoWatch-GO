FROM node:20-bullseye

WORKDIR /app

# Only copy package files first
COPY package*.json ./

# Install dependencies inside container
RUN npm install

# Copy the rest of the source code
COPY . .

# Build the app
RUN npm run build

# Expose port for adapter-node
EXPOSE 3000
CMD ["node", "build/index.js"]