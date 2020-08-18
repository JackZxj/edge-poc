# EDGE POC

## env

kubernetes version: v1.17.0
kubeedge version: v1.3.1

raspberry-zero w1.1 with RaspbianOS jessie

## prepare os

``` BASH
# download OS image
wget http://downloads.raspberrypi.org/raspbian_lite/images/raspbian_lite-2017-07-05/2017-07-05-raspbian-jessie-lite.zip
unzip 2017-07-05-raspbian-jessie-lite.zip

# download image tools: Win32DiskImager
# you can get it from here: https://sourceforge.net/projects/win32diskimager/
# install it and then Write the image to your SD card
```

### cross-compiling kernel

Ubuntu 16.04

* 2core 64bit CPU
* 4GB RAM

``` BASH
# change to root
su - root
# Install required dependencies
apt install git bc bison flex libssl-dev make libc6-dev libncurses5-dev

# Install toolchain
# mirror: https://gitee.com/x_jackzhang/raspberrypi-tools.git
git clone https://github.com/raspberrypi/tools ~/tools

# add it to PATH
echo PATH=\$PATH:~/tools/arm-bcm2708/gcc-linaro-arm-linux-gnueabihf-raspbian-x64/bin >> ~/.bashrc
source ~/.bashrc
# Get sources
git clone --depth=1 --branch rpi-4.9.y https://github.com/raspberrypi/linux

# After a long long wait
cd linux
KERNEL=kernel
# add CONFIG_CGROUP_PIDS=y to it bcmrpi_defconfig
vi arch/arm/configs/bcmrpi_defconfig
# prepare .config
make ARCH=arm CROSS_COMPILE=arm-linux-gnueabihf- bcmrpi_defconfig
# build the kernel
# To speed up compilation on multiprocessor systems, and get some improvement on single processor ones, use -j n, where n is the number of processors * 2. Alternatively, feel free to experiment and see what works!
make ARCH=arm CROSS_COMPILE=arm-linux-gnueabihf- zImage modules dtbs

# Plug your SD card to Ubuntu, you may get something such as these
lsblk
# NAME   MAJ:MIN RM  SIZE RO TYPE MOUNTPOINT
# sdb      8:16   1 58.2G  0 disk 
# ├─sdb2   8:18   1  1.6G  0 part 
# └─sdb1   8:17   1 41.7M  0 part
mkdir -p /mnt/fat32
mkdir /mnt/ext4
mount /dev/sdb1 /mnt/fat32
mount /dev/sdb2 /mnt/ext4

# Next, install the modules
make ARCH=arm CROSS_COMPILE=arm-linux-gnueabihf- INSTALL_MOD_PATH=/mnt/ext4 modules_install

# Finally, copy the kernel and Device Tree blobs onto the SD card, making sure to back up your old kernel:
cp /mnt/fat32/$KERNEL.img /mnt/fat32/$KERNEL-backup.img
cp arch/arm/boot/zImage /mnt/fat32/$KERNEL.img
cp arch/arm/boot/dts/*.dtb /mnt/fat32/
cp arch/arm/boot/dts/overlays/*.dtb* /mnt/fat32/overlays/
cp arch/arm/boot/dts/overlays/README /mnt/fat32/overlays/

# enable ssh onboot
touch /mnt/fat32/ssh

# enable wifi onboot
touch /mnt/fat32/wpa_supplicant.conf
cat > /mnt/fat32/wpa_supplicant.conf <<EOF
country=CN
ctrl_interface=DIR=/var/run/wpa_supplicant GROUP=netdev
update_config=1

network={
    ssid="WIFI NAME"
    psk="WIFI PASSWORD"
    key_mgmt=WPA-PSK
}
EOF

# config cgroup
sed -i  's/rootwait/& cgroup_enable=cpuset cgroup_enable=memory cgroup_memory=1/' /mnt/fat32/cmdline.txt

umount /mnt/fat32
umount /mnt/ext4

# Plug the SD card to raspberrypi zero w
```

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

## prepare edge

``` BASH
# mosquitto server
apt-get install mosquitto
# client for test mosquitto
apt-get install mosquitto-clients
# update mqtt config
vi /etc/mosquitto/conf.d/port.conf
port 1883
listener 1884
# restart mqtt
/etc/init.d/mosquitto restart
# or start manually:
# mosquitto -c /etc/mosquitto/mosquitto.conf -d

# test
# terminal 1: for subscribe
mosquitto_sub -t hello
# terminal 2: for publish
mosquitto_pub -t hello -h localhost -m "hello world!" 
```

## kubeedge for test

``` BASH

#################### setup cloud side ####################

# prepare kubeedge cloudcore config file
./cloudcore --minconfig > cloudcore.yaml

# run cloudcore(fg)
./cloudcore --config cloudcore.yaml
# run cloudcore(bg)
nohup ./cloudcore --config cloudcore.yaml > cloudcore.log 2>&1 &
#################### setup edge side ####################

# prepare kubeedge edgecore config file
./edgecore --minconfig > edgecore.yaml

# get token form k8s master
kubectl get secret -nkubeedge tokensecret -o=jsonpath='{.data.tokendata}' | base64 -d

# insert token into edgecore config
sed -i -e "s|token: .*|token: ${token}|g" edgecore.yaml

# run edgecore(fg)
./edgecore --config edgecore.yaml
# run edgecore(bg)
nohup ./edgecore --config edgecore.yaml > edgecore.log 2>&1 &
```

## kubeedge for production 

...


## run demo

``` BASH
# run on edge site
cd temperature-and-humidity
docker build -f Dockerfile-DHT11 -t edge-temperature-and-humidity:v1 .

# run on cloud site
kubectl apply ./temperature-and-humidity/crds/devicemodel.yaml
kubectl apply ./temperature-and-humidity/crds/device.yaml
kubectl apply ./temperature-and-humidity/deployment.yaml
```
