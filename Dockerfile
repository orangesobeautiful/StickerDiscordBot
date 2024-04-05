FROM docker.io/library/node:14.21.3 AS frontend-deps
WORKDIR /frontend/
RUN npm i -g @quasar/cli@2.3.0
COPY ./web/frontend/package.json ./web/frontend/yarn.lock /frontend/
RUN npm install

FROM frontend-deps AS build-frontend
WORKDIR /frontend/
COPY ./web/frontend/ /frontend/
RUN quasar build -m spa

FROM docker.io/library/golang:1.21.4-alpine as build-backend
WORKDIR /backend/
RUN go version
COPY ./backend/go.mod ./backend/go.sum /backend/
RUN go mod download
COPY ./backend/ /backend/
RUN CGO_ENABLED=0 go build -ldflags="-s" -o backend-server
COPY --from=build-frontend /frontend/dist/spa/ /backend/frontend-dist/frontend-static/original/

FROM docker.io/library/python:3.10-slim as base
FROM base as bot_builder
COPY requirements.txt .
RUN pip install --user -r requirements.txt
FROM base
# install bot dependencies
RUN apt-get update && \
    apt-get install --no-install-recommends -y libopus0=1.3.1-3 supervisor=4.2.5-1 && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*
# copy bot dependencies
COPY --from=bot_builder /root/.local /root/.local
COPY ./supervisord.conf /etc/supervisor/conf.d/supervisord.conf
ENV PATH=/root/.local/bin:$PATH
WORKDIR /app/
COPY --from=build-frontend /frontend/dist/spa/ /app/public/
COPY --from=build-backend /backend/backend-server /app/
COPY ./bot/ /app/bot/
COPY ./CommonFunction/ /app/CommonFunction/
COPY ./Controller/ /app/Controller/
COPY ./Database/ /app/Database/
COPY ./runBot.py /app/

CMD [ "supervisord" ]