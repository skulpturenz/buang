FROM golang:alpine AS build

WORKDIR /build

COPY . .

# TODO: BUMP
ARG VERSION="unknown"
RUN go build -ldflags "-X skulpture/buang/constants/envs.BUANG_VERSION=$VERSION"

FROM alpine

RUN apk update && \
	apk add --no-cache docker docker-cli busybox-openrc && \
	mkdir -p /etc/docker && \
	echo "{ 'features': { 'buildkit': true } }" > /etc/docker/daemon.json && \
	rc-service docker start && \
	rc-update add docker boot

COPY --from=build /build/buang /app/buang

ENTRYPOINT [ "/app/buang" ]
