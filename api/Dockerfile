FROM golang:alpine AS build

WORKDIR /build

COPY . .

# TODO: BUMP
ARG VERSION="unknown"
RUN go build -ldflags "-X skulpture/buang/constants/envs.BUANG_VERSION=$VERSION"

FROM alpine

RUN apk update && \
	mkdir -p /etc/docker && \
	echo "{ \"features\": { \"buildkit\": true } }" > /etc/docker/daemon.json && \
	apk add --no-cache docker docker-cli busybox-openrc && \
	# TODO: rc-service returns a non zero error, too lazy to figure out why
	(rc-service --quiet docker restart || true) && \
	(rc-update --quiet add docker boot || true)

COPY --from=build /build/buang /app/buang

ENTRYPOINT [ "/app/buang" ]
