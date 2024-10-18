FROM debian:bookworm-slim

COPY centralized-relay /usr/bin/centralized-relay

RUN useradd -ms /bin/bash relayer
WORKDIR /home/relayer

USER relayer

ENTRYPOINT [ "/usr/bin/centralized-relay" ]