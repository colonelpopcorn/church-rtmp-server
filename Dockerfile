FROM alpine:3.14

COPY / /tmp/project
WORKDIR /tmp/project/provisioning
RUN \
apk add --no-cache --virtual .ansible ansible \
&& ansible-playbook --become --connection=local -i 127.0.0.1, bootstrap.yml \
&& apk del .ansible

EXPOSE 80 443 1935

ENTRYPOINT [ "/usr/bin/systemctl", "start" ]
CMD [ "nginx", "streaming-server" ]