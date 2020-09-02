FROM arm32v6/golang:1.12-alpine3.10

RUN echo http://mirrors.ustc.edu.cn/alpine/v3.10/main > /etc/apk/repositories; \
    echo http://mirrors.ustc.edu.cn/alpine/v3.10/community >> /etc/apk/repositories; \
    apk add gcc libc-dev --no-cache

COPY . $GOPATH/src/RGB

# NOTE: Pwm is possible only for pins 12, 13, 18, 19. Pin mode: BCM
# https://github.com/stianeikeland/go-rpio/blob/v4.2.0/rpio.go#L238
ENV GO111MODULE=off RPIN=12 GPIN=13 BPIN=19 DEVICE_NAME=rgb-light-device

WORKDIR $GOPATH/src/RGB/src

RUN go build -v -o /go/bin/rgb-light && \
    chmod u+x /go/bin/rgb-light

WORKDIR /go/bin/

CMD ["rgb-light", "-logtostderr=true"]