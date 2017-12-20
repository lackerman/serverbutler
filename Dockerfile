FROM golang:1.9.2 as builder
WORKDIR /go/src/github.com/lackerman/serverbutler
COPY . ./

RUN go get -d -v golang.org/x/net/html \
	&& go get -u github.com/jteeuwen/go-bindata/... \
	&& go get ./... \
	&& CGO_ENABLED=0 GOOS=linux ./build.sh

FROM linuxserver/transmission
LABEL AUTHOR=lackerman

# install packages
RUN apk add --no-cache openvpn bash \
	&& mkdir -p /dev/net \
	&& mknod /dev/net/tun c 10 200 \
	&& chmod 600 /dev/net/tun

COPY --from=builder /go/src/github.com/lackerman/serverbutler/bin/serverbutler /bin

COPY root/ /

WORKDIR /go/src/github.com/lackerman/serverbutler/bin

# ports and volumes
EXPOSE 8080
