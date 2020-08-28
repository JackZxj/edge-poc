## run demo

``` BASH
# run on edge side
docker build -t edge-rgb-light:v1 .

# run on cloud side
kubectl apply -f ./crds/devicemodel.yaml
kubectl apply -f ./crds/device.yaml

# update before apply
kubectl apply -f ./deployment.yaml

# get info
kubectl get device rgb-light-device -oyaml -w
```

## issue

there are some problem with `go-rpio` in pwm mode.

1. `go-rpio` was a hardware method for pwm, only support pins: 12, 13, 18, 19
2. `go-rpio` didn't work as expected. I use pins 12, 13, 19 to control, but pin 13 didn't work. This can be verified in `demo` .

https://github.com/stianeikeland/go-rpio/issues/20
