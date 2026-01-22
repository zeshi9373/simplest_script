FROM alpine

WORKDIR /app

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories

RUN apk add tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

COPY simplest_script /app/simplest_script
COPY etc/release.yaml /app/release.yaml

CMD ["./simplest_script", "-f", "release.yaml"]