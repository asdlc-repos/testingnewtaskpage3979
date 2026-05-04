#!/bin/sh
# Strip trailing slashes — OpenChoreo may inject "http://host:8080/"
LEAVE_SERVICE_URL="${LEAVE_SERVICE_URL%/}"
USER_SERVICE_URL="${USER_SERVICE_URL%/}"

# Default to localhost if not set
LEAVE_SERVICE_URL="${LEAVE_SERVICE_URL:-http://localhost:9090}"
USER_SERVICE_URL="${USER_SERVICE_URL:-http://localhost:9091}"

# Substitute only our variables, protecting nginx's own $variables
envsubst '$LEAVE_SERVICE_URL $USER_SERVICE_URL' \
  < /etc/nginx/conf.d/default.conf.template \
  > /etc/nginx/conf.d/default.conf

cat <<EOF > /usr/share/nginx/html/env.js
window.RUNTIME_LEAVE_SERVICE_URL = "/api/leave";
window.RUNTIME_USER_SERVICE_URL = "/api/user";
EOF

exec "$@"
