FROM ubuntu:18.04
WORKDIR /
COPY --chown=root:root ./bin/elaina /
USER root:root
ENTRYPOINT ["tail","-f","/dev/null"]