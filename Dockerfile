#use node alpine v14.15.1 - I have 14.15.2 on my Ubuntu 18.04LTS
FROM node:14.15.1-alpine3.10

COPY . /app

RUN npm install -g @angular/cli

