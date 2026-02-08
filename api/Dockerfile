FROM golang:alpine AS build

WORKDIR /build

COPY . .

# TODO: BUMP
ARG VERSION="unknown"
RUN go build -ldflags "-X skulpture/buang/constants/envs/envs.BUANG_VERSION=$VERSION"

FROM alpine

RUN apk add --no-cache docker-cli

COPY --from=build /build/buang /app/buang

ENTRYPOINT [ "/app/buang" ]
