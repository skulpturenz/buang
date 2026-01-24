# TODO: BUMP
ARG VERSION="0.0.1"

FROM golang:alpine AS build

WORKDIR /build

COPY . .

RUN go build -ldflags "-X skulpture/buang/constants/envs.BUANG_VERSION=$VERSION"

FROM alpine

COPY --from=build /build/buang /buang

ENTRYPOINT [ "/buang" ]
