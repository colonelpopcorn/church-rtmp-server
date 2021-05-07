FROM alpine:3.13.5

COPY / /tmp/project
WORKDIR /tmp/project/provisioning
RUN \
apk add --no-cache --virtual .ansible ansible \
&& apk add --no-cache --virtual .python3 python3 \
&& wget -O /usr/bin/systemctl https://raw.githubusercontent.com/gdraheim/docker-systemctl-replacement/master/files/docker/systemctl3.py \
&& chmod +x /usr/bin/systemctl \
&& ansible-playbook --become --connection=local -i 127.0.0.1, bootstrap.yml \
&& rm -rf /root/* \
&& rm -rf /tmp/project/* \
&& rm /usr/local/bin/go /usr/local/bin/gofmt \
&& rm -rf /usr/local/go \
&& apk del .ansible .build-deps

EXPOSE 80 443 1935

ENTRYPOINT [ "/usr/bin/systemctl", "start" ]
CMD [ "nginx", "streaming-server" ]