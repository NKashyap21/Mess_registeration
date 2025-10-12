FROM node:latest

RUN npm i -g pnpm

COPY package.json ./
COPY pnpm-lock.yaml ./
COPY build/ ./build/

RUN pnpm install --frozen-lockfile --prod

ENV PORT=3000
ENV ORIGIN=https://mess.iith.dev
ENV NODE_ENV='production'

EXPOSE 3000

CMD ["node",  "build"]
