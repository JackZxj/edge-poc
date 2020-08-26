FROM arm32v6/golang:1.12-alpine3.10

RUN echo http://mirrors.ustc.edu.cn/alpine/v3.10/main > /etc/apk/repositories; \
    echo http://mirrors.ustc.edu.cn/alpine/v3.10/community >> /etc/apk/repositories; \
    apk add gcc libc-dev --no-cache

COPY . $GOPATH/src/RGB

ENV GO111MODULE=off RPIN=21 GPIN=20 BPIN=16

WORKDIR $GOPATH/src/RGB/src

RUN go build -v -o /go/bin/rgb-light && \
    chmod u+x /go/bin/rgb-light

WORKDIR /go/bin/

CMD ["rgb-light", "-logtostderr=true"]