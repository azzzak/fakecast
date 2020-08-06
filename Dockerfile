FROM golang:1.14-alpine AS backend
ARG APP_VER
WORKDIR /fakecast
COPY backend ./
RUN apk --no-cache --update add \
  musl-dev \
  gcc
RUN go mod tidy && \
  go mod download && \
  go build -o app -ldflags "-X main.version=$APP_VER -s -w" ./

FROM node:14.5-alpine AS frontend
WORKDIR /front
COPY frontend/package.json frontend/package-lock.json ./
COPY frontend/src ./src
COPY frontend/public ./public
RUN npm ci --silent
RUN npm install react-scripts@3.4.1 -g --silent
RUN npm run build

FROM alpine:3.12
RUN apk --no-cache --update add libcap && \
  addgroup -g 1000 -S fakecast && \
  adduser -u 1000 -S fakecast -G fakecast && \
  mkdir /fakecast && \
  chown -R fakecast:fakecast /fakecast
WORKDIR /home/fakecast
COPY --from=backend --chown=fakecast:fakecast /fakecast/app ./
COPY --from=frontend --chown=fakecast:fakecast /front/build ./front
RUN setcap 'cap_net_bind_service=+ep' ./app
USER fakecast

EXPOSE 80
ENTRYPOINT ["./app"]