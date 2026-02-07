# Stage 1: Build backend
FROM golang:1.24-bookworm AS backend-builder
WORKDIR /build
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ .
RUN CGO_ENABLED=0 go build -o virtpanel ./cmd/main.go

# Stage 2: Build frontend
FROM node:20-slim AS frontend-builder
RUN npm install -g pnpm
WORKDIR /build
COPY frontend/package.json frontend/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile
COPY frontend/ .
RUN pnpm build

# Stage 3: Runtime
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    qemu-kvm qemu-utils libvirt-daemon-system virtinst \
    dnsmasq-base iptables iproute2 nginx \
    && rm -rf /var/lib/apt/lists/*

# Copy backend binary
COPY --from=backend-builder /build/virtpanel /usr/local/bin/virtpanel

# Copy frontend dist
COPY --from=frontend-builder /build/dist /var/www/virtpanel

# Nginx config
COPY docker/nginx.conf /etc/nginx/sites-available/default

# Startup script
COPY docker/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 80

ENTRYPOINT ["/entrypoint.sh"]
