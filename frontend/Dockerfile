# base image
FROM node:22-alpine AS base
LABEL maintainer="takatost@gmail.com"

# if you located in China, you can use aliyun mirror to speed up
# RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# if you located in China, you can use taobao registry to speed up
# RUN npm config set registry https://registry.npmmirror.com

RUN apk add --no-cache tzdata
RUN corepack enable
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
ARG NEXT_PUBLIC_BASE_PATH=""
ENV NEXT_PUBLIC_BASE_PATH="$NEXT_PUBLIC_BASE_PATH"


# install packages
FROM base AS packages

WORKDIR /app/web

COPY package.json pnpm-lock.yaml /app/web/

# Use packageManager from package.json
RUN corepack install

RUN pnpm install --frozen-lockfile

# build resources
FROM base AS builder
WORKDIR /app/web
COPY --from=packages /app/web/ .
COPY . .

ENV NODE_OPTIONS="--max-old-space-size=4096"
RUN pnpm build


# production stage
FROM base AS production

ENV NODE_ENV=production
ENV EDITION=SELF_HOSTED
ENV DEPLOY_ENV=PRODUCTION
ENV CONSOLE_API_URL=http://127.0.0.1:5001
ENV APP_API_URL=http://127.0.0.1:5001
ENV MARKETPLACE_API_URL=https://marketplace.gofy.ai
ENV MARKETPLACE_URL=https://marketplace.gofy.ai
ENV PORT=3000
ENV NEXT_TELEMETRY_DISABLED=1

# set timezone
ENV TZ=UTC
RUN ln -s /usr/share/zoneinfo/${TZ} /etc/localtime \
    && echo ${TZ} > /etc/timezone

# Create non-root user
ARG gofy_uid=1001
RUN addgroup -S -g ${gofy_uid} gofy && \
    adduser -S -u ${gofy_uid} -G gofy -s /bin/ash -h /home/gofy gofy && \
    mkdir /app && \
    chown -R gofy:gofy /app


WORKDIR /app/web

COPY --from=builder --chown=gofy:gofy /app/web/public ./public
COPY --from=builder --chown=gofy:gofy /app/web/.next/standalone ./
COPY --from=builder --chown=gofy:gofy /app/web/.next/static ./.next/static

COPY --chown=gofy:gofy --chmod=755 docker/entrypoint.sh ./entrypoint.sh

ARG COMMIT_SHA
ENV COMMIT_SHA=${COMMIT_SHA}

USER gofy
EXPOSE 3000
ENTRYPOINT ["/bin/sh", "./entrypoint.sh"]
