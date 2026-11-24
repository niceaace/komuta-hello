#!/bin/sh
set -e
PORT="${PORT:-8080}"
if [ -z "$VLESS_UUID" ]; then echo "VLESS_UUID not set"; exit 1; fi
cat > /tmp/config.json <<EOF
{
  "log": { "level": "info" },
  "inbounds": [
    {
      "type": "vless",
      "tag": "vless-in",
      "listen": "0.0.0.0",
      "listen_port": $PORT,
      "users": [ { "uuid": "$VLESS_UUID" } ],
      "transport": { "type": "ws", "path": "/raycdn" }
    }
  ],
  "outbounds": [ { "type": "direct", "tag": "direct" } ]
}
EOF
exec sing-box run -c /tmp/config.json











