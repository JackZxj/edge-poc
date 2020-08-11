# EDGE POC

## env

kubernetes version: v1.17.0
kubeedge version: v1.3.1

raspberry-zero w1.1 with RaspbianOS

## prepare k8s

``` BASH
# k8s init
kubeadm init --pod-network-cidr=10.244.0.0/16 --image-repository=registry.aliyuncs.com/google_containers --kubernetes-version=v1.17.0

mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config

# set flannel network
# official: https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel.yml
kubectl apply -f kube-flannel.yml

# for kubeedge cloudcore
kubectl apply -f ./CRDs/devices/devices_v1alpha1_device.yaml
kubectl apply -f ./CRDs/devices/devices_v1alpha1_devicemodel.yaml
kubectl apply -f ./CRDs/reliablesyncs/cluster_objectsync_v1alpha1.yaml
kubectl apply -f ./CRDs/reliablesyncs/objectsync_v1alpha1.yaml
```

## kubeedge for test

``` BASH

#################### setup cloud side ####################

# prepare kubeedge cloudcore config file
cloudcore --minconfig > cloudcore.yaml

# run cloudcore
cloudcore --config cloudcore.yaml

#################### setup edge side ####################

# prepare kubeedge edgecore config file
edgecore --minconfig > edgecore.yaml

# get token form k8s master
kubectl get secret -nkubeedge tokensecret -o=jsonpath='{.data.tokendata}' | base64 -d

# insert token into edgecore config
sed -i -e "s|token: .*|token: ${token}|g" edgecore.yaml

# run
edgecore --config edgecore.yaml
```

## kubeedge for production 

...
