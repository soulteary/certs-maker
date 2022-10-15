FROM golang:1.19.0-alpine3.16 as builder
RUN echo '' > /etc/apk/repositories && \
    echo "https://mirror.tuna.tsinghua.edu.cn/alpine/v3.16/main"         >> /etc/apk/repositories && \
    echo "https://mirror.tuna.tsinghua.edu.cn/alpine/v3.16/community"    >> /etc/apk/repositories && \
    echo "Asia/Shanghai" > /etc/timezone
RUN apk add upx
WORKDIR /build
ENV CGO_ENABLED=0
COPY ./app .
RUN go build -ldflags "-w -s"
RUN upx -9 -o cert-maker.minify cert-maker && mv cert-maker.minify cert-maker

FROM alpine:3.16
RUN echo '' > /etc/apk/repositories && \
    echo "https://mirror.tuna.tsinghua.edu.cn/alpine/v3.16/main"         >> /etc/apk/repositories && \
    echo "https://mirror.tuna.tsinghua.edu.cn/alpine/v3.16/community"    >> /etc/apk/repositories && \
    echo "Asia/Shanghai" > /etc/timezone
RUN apk add openssl && rm -rf /var/cache/apk/*
WORKDIR /
COPY --from=Builder /build/cert-maker /bin/cert-maker
COPY --from=Builder /build/entrypoint.sh /entrypoint.sh
CMD ["/entrypoint.sh"]