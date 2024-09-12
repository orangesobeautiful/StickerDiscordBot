# Create flutter environment
FROM --platform=$BUILDPLATFORM dart:3.5.2 AS dart_base

RUN curl -O https://storage.googleapis.com/flutter_infra_release/releases/stable/linux/flutter_linux_3.24.2-stable.tar.xz

RUN mkdir -p /development

RUN apt-get update && \
    apt-get install --no-install-recommends -y xz-utils && \
    rm -rf /var/lib/apt/lists/*

RUN tar -xf flutter_linux_3.24.2-stable.tar.xz -C /development && \
    rm -rf flutter_linux_3.24.2-stable.tar.xzrm -rf flutter_linux_3.24.2-stable.tar.xz

ENV PATH="/development/flutter/bin:${PATH}"

RUN git config --global --add safe.directory /development/flutter

# build frontend
FROM --platform=$BUILDPLATFORM dart_base AS flutter_builder

COPY ./frontend /app

WORKDIR /app

RUN flutter pub get

RUN flutter pub run build_runner build --delete-conflicting-outputs

RUN flutter build web

# build backend
FROM --platform=$BUILDPLATFORM golang:1.23.1-bookworm AS go_builder

WORKDIR /app

COPY ./backend/go.mod ./backend/go.sum ./

RUN go mod download

COPY ./backend .

RUN go generate ./...

ARG TARGETOS TARGETARCH

RUN GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 go build -ldflags="-s" -o /app/server .

# runtime image
FROM alpine:3.20.3

WORKDIR /app

COPY --from=go_builder /app/server .
COPY ./backend/migrations ./migrations
COPY --from=flutter_builder /app/build/web ./frontend-web

CMD ["/app/server"]
