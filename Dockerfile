FROM alpine:3.14

COPY / /tmp/project
WORKDIR /tmp/project/provisioning
RUN \
    apk add --no-cache --virtual .ansible ansible \
    && ansible-playbook --become --connection=local -i 127.0.0.1, bootstrap.yml \
    && apk del .ansible \
    && rm -rf /tmp/project/* 

EXPOSE 80 443 1935 9001

ENTRYPOINT [ "/usr/local/bin/supervisord/supervisord_static", "-c"]
CMD [ "/etc/supervisor.conf" ]