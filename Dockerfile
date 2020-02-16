FROM alpine
ADD built/kubano /bin/
ENTRYPOINT /bin/kubano
