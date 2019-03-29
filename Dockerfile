FROM alpine
ADD slack /bin/
RUN apk -Uuv add ca-certificates
ENTRYPOINT /bin/slack