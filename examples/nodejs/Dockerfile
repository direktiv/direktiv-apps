FROM node:12-alpine

WORKDIR /app

COPY package*.json ./
COPY index.js ./

RUN npm ci --only=production
EXPOSE 8080

CMD ["node", "index.js"]