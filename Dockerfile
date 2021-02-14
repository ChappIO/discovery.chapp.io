FROM node:15 as build

RUN mkdir /build
WORKDIR /build

ADD package.json .
ADD yarn.lock .
RUN yarn --frozen-lockfile
ADD *.json .
ADD ./src ./src
RUN yarn build

FROM node:15

ENV NODE_ENV=production

RUN mkdir /app
WORKDIR /app

ADD package.json .
ADD yarn.lock .
RUN yarn --frozen-lockfile --production
COPY --from=build /build/dist ./dist

EXPOSE 3000
ENTRYPOINT ["node", "/app/dist/main.js"]
