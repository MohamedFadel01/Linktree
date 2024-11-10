FROM node:20-alpine AS builder

WORKDIR /app

COPY package*.json ./

RUN npm ci

COPY . .

ARG VITE_API_URL
ENV VITE_API_URL=${VITE_API_URL}

RUN npm run build

FROM node:20-alpine

WORKDIR /app

COPY package*.json ./

RUN npm ci

COPY --from=builder /app/dist ./dist

COPY vite.config.js ./

EXPOSE 3000

ENV NODE_ENV=production

CMD ["npx", "vite", "preview", "--host", "0.0.0.0", "--port", "3000"]
