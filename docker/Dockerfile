FROM golang:1.20.0-alpine3.16 as builder
RUN echo '' > /etc/apk/repositories && \
    echo "https://mirror.tuna.tsinghua.edu.cn/alpine/v3.16/main"         >> /etc/apk/repositories && \
    echo "https://mirror.tuna.tsinghua.edu.cn/alpine/v3.16/community"    >> /etc/apk/repositories && \
    echo "Asia/Shanghai" > /etc/timezone
RUN apk add upx
WORKDIR /build
ENV CGO_ENABLED=0
COPY . .
RUN go build -ldflags "-w -s"
RUN upx -9 -o certs-maker.minify certs-maker && mv certs-maker.minify certs-maker

FROM alpine:3.16
RUN echo '' > /etc/apk/repositories && \
    echo "https://mirror.tuna.tsinghua.edu.cn/alpine/v3.16/main"         >> /etc/apk/repositories && \
    echo "https://mirror.tuna.tsinghua.edu.cn/alpine/v3.16/community"    >> /etc/apk/repositories && \
    echo "Asia/Shanghai" > /etc/timezone
RUN apk add openssl && rm -rf /var/cache/apk/*
WORKDIR /
COPY --from=Builder /build/certs-maker /bin/certs-maker
ENTRYPOINT ["certs-maker"]