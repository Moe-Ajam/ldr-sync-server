FROM debian:stable-slim

# COPY source destination
COPY video-sync-server /bin/video-sync-server

CMD ["/bin/video-sync-server"]
