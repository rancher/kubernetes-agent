FROM debian:jessie
RUN apt-get update -y && \
    apt-get -yy -q install ca-certificates
COPY ./kubernetes-agent /usr/bin/kubernetes-agent
CMD ["kubernetes-agent"]
