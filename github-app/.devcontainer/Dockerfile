FROM node:lts AS base

RUN apt update && apt upgrade -y && \
    apt install vim bash-completion -y && \
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.3/migrate.linux-$(dpkg --print-architecture).deb -o migrate.deb && \
    dpkg -i migrate.deb
RUN git config --global pager.branch false
RUN echo "if [ -f /etc/bash_completion ]; then\n\t. /etc/bash_completion\nfi" >> /etc/bash.bashrc

FROM base

RUN bash -c "corepack enable pnpm && pnpm config set store-dir /home/node/.local/share/pnpm/store"
