FROM alpine:3.15

WORKDIR /opt/postgres-init

COPY init-scripts init-scripts
COPY run.sh run.sh
RUN chmod +x run.sh

ENTRYPOINT ["sh", "-c", "/opt/postgres-init/run.sh"]
