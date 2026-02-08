FROM golang:alpine AS build

# TODO: BUMP
ARG VERSION="unknown"

WORKDIR /build

COPY . .

RUN go build -ldflags "-X skulpture/buang/constants/envs.BUANG_VERSION=$VERSION"

FROM alpine

RUN apk add --no-cache docker-cli

COPY --from=build /build/buang /app/buang

ENTRYPOINT [ "/app/buang" ]
