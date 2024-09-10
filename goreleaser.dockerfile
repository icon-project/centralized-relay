FROM scratch

COPY centralized-relay .

ENTRYPOINT [ "centralized-relay" ]