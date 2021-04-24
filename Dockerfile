FROM ubuntu:20.04

RUN apt update && apt install ansible -y
COPY / /tmp/project

RUN wget -O /usr/bin/systemctl https://raw.githubusercontent.com/gdraheim/docker-systemctl-replacement/master/files/docker/systemctl3.py
RUN chmod +x /usr/bin/systemctl

WORKDIR /tmp/project/provisioning
RUN ansible-playbook --become --connection=local -i 127.0.0.1, bootstrap.yml

EXPOSE 80 443 1935

ENTRYPOINT [ "/usr/bin/systemctl", "start" ]
CMD [ "nginx", "streaming-server" ]