FROM arm32v6/golang:1.12-alpine3.10

COPY . $GOPATH/src/github.com/kubeedge/examples/led-raspberrypi

ENV GO111MODULE=off

WORKDIR $GOPATH/src/github.com/kubeedge/examples/led-raspberrypi

RUN go build -v -o light_mapper

CMD ["./light_mapper", "-logtostderr=true"]
