FROM alpine
ADD kubano /bin/
RUN apk -Uuv add ca-certificates
ENTRYPOINT /bin/kubano
