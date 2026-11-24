FROM node:20-alpine
RUN npm install -g npm@latest
WORKDIR /app
COPY . .
EXPOSE 8080
CMD ["node", "server.js"]






