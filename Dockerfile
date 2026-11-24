FROM node:20-alpine
RUN npm install -g npm@10
WORKDIR /app
COPY . .
EXPOSE 8080
CMD ["node", "server.js"]