FROM alpine
COPY built/kubano /bin/
ENTRYPOINT /bin/kubano
