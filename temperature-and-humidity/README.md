## run demo

``` BASH
# run on edge side
docker build -f Dockerfile-DHT11 -t edge-temperature-and-humidity:v1 .

# run on cloud side
kubectl apply -f ./crds/devicemodel.yaml
kubectl apply -f ./crds/device.yaml

# update before apply
kubectl apply -f ./deployment.yaml

# get info
kubectl get device temperature-and-humidity -oyaml -w
```

``` yaml
apiVersion: devices.kubeedge.io/v1alpha1
kind: Device
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"devices.kubeedge.io/v1alpha1","kind":"Device","metadata":{"annotations":{},"labels":{"description":"temperature-and-humidity"},"name":"temperature-and-humidity","namespace":"default"},"spec":{"deviceModelRef":{"name":"tmp-rh-model"},"nodeSelector":{"nodeSelectorTerms":[{"matchExpressions":[{"key":"kubernetes.io/hostname","operator":"In","values":["edge-zero02"]}]}]}},"status":{"twins":[{"desired":{"metadata":{"type":"string"},"value":""},"propertyName":"temperature-status"},{"desired":{"metadata":{"type":"string"},"value":""},"propertyName":"humidity-status"}]}}
  creationTimestamp: "2020-08-18T07:00:07Z"
  generation: 16139
  labels:
    description: temperature-and-humidity
  name: temperature-and-humidity
  namespace: default
  resourceVersion: "852578"
  selfLink: /apis/devices.kubeedge.io/v1alpha1/namespaces/default/devices/temperature-and-humidity
  uid: b1c0bb07-d7d3-4694-9168-897ebeca9841
spec:
  deviceModelRef:
    name: tmp-rh-model
  nodeSelector:
    nodeSelectorTerms:
    - matchExpressions:
      - key: kubernetes.io/hostname
        operator: In
        values:
        - edge-zero02
status:
  twins:
  - desired:
      metadata:
        type: string
      value: ""
    propertyName: temperature-status
    reported:
      metadata:
        timestamp: "1597885303105"
        type: string
      value: 26C
  - desired:
      metadata:
        type: string
      value: ""
    propertyName: humidity-status
    reported:
      metadata:
        timestamp: "1597885303107"
        type: string
      value: 79#
```