FROM ghcr.io/sagernet/sing-box:latest
COPY start.sh /start.sh
ENTRYPOINT ["/bin/sh", "/start.sh"]



