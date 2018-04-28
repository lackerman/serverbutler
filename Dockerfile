FROM alpine
LABEL AUTHOR=lackerman

ARG config_dir=/tmp
ENV config_dir=$config_dir

WORKDIR /opt/sb

RUN apk add --no-cache ca-certificates

EXPOSE 3000

CMD ["./serverbutler"]

COPY bin/serverbutler .