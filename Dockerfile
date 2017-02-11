FROM alpine:latest

WORKDIR "/opt"

ADD .docker_build/udoit /opt/bin/udoit

CMD ["/opt/bin/udoit"]
