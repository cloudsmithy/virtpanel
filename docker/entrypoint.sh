#!/bin/bash
set -e

# Start libvirt daemons
virtlogd -d
libvirtd -d
sleep 1

# Activate default NAT network
virsh net-start default 2>/dev/null || true
virsh net-autostart default 2>/dev/null || true

# Restore port forward rules
# (backend does this on startup)

# Start nginx
nginx

# Start backend (foreground)
exec /usr/local/bin/virtpanel
