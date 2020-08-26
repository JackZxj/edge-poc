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