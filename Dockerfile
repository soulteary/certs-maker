FROM node:15.7.0-alpine3.12
RUN apk add openssl

WORKDIR /app
COPY app/index.js app/package.json app/package-lock.json ./
RUN npm install 

WORKDIR /
COPY app/entrypoint.sh ./

ENTRYPOINT ["/entrypoint.sh"]
