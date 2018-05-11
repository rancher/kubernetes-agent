FROM rancher/agent-base:v0.3.0
RUN apt-get update -y && \
    apt-get -yy -q install ca-certificates wget unzip
ENV SSL_SCRIPT_COMMIT 98660ada3d800f653fc1f105771b5173f9d1a019
RUN wget -O /usr/bin/update-rancher-ssl https://raw.githubusercontent.com/rancher/rancher/${SSL_SCRIPT_COMMIT}/server/bin/update-rancher-ssl && \
    chmod +x /usr/bin/update-rancher-ssl && \
    curl -sLf https://github.com/leodotcloud/loglevel/releases/download/v0.1/loglevel-amd64-v0.1.tar.gz | tar xvzf - -C /usr/bin
COPY ./kubernetes-agent /usr/bin/kubernetes-agent
COPY ./entry.sh /usr/bin/entry.sh
ENTRYPOINT ["/usr/bin/entry.sh"]
CMD ["kubernetes-agent"]
