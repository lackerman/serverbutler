# build stage
FROM golang:alpine AS build-env
RUN apk --no-cache add build-base git bzr mercurial gcc make
ENV GO111MODULE=on

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN make build

# final stage
FROM alpine
LABEL AUTHOR=lackerman

WORKDIR /opt
RUN apk add --no-cache ca-certificates curl

ENTRYPOINT ["./app"]

COPY --from=build-env /app/bin/app ./app
COPY --from=build-env /app/templates/* ./templates/