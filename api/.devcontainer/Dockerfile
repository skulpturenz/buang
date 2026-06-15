FROM golang:1-alpine AS base

RUN apk add git bash shadow gnupg curl python3 py3-pip bash-completion docker-cli && echo '/bin/bash' | chsh `whoami`
RUN git config --global pager.branch false

COPY ./kryptos.sh /usr/local/bin/kryptos
RUN chmod u+x /usr/local/bin/kryptos

FROM base

RUN apk update && \
	apk add --no-cache build-base libc6-compat git-lfs ca-certificates openssl nodejs npm git-cliff && \
	npm i -g @google/gemini-cli && \
	git lfs install && \
	update-ca-certificates

RUN bash -c "curl https://sdk.cloud.google.com > /tmp/install-gcloud.sh" && bash /tmp/install-gcloud.sh --disable-prompts
ENV PATH="${PATH}:/root/google-cloud-sdk/bin"
