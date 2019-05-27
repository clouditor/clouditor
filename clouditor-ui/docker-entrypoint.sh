#/bin/bash
set -a

if [ -n "$API_URL" ]; then
  echo {\"apiUrl\": \"${API_URL}\"} > /usr/share/nginx/html/assets/config.json
fi

if [ "$1"  = 'nginx' ]; then
  exec nginx -g "daemon off;"
fi

exec "$@"
