FROM node

WORKDIR /app/frontend
COPY package.json /app/frontend
COPY yarn.lock /app/frontend

RUN yarn install

EXPOSE 3000
CMD ["yarn", "start"]