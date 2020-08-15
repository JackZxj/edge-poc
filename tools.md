# Tools for raspberry-pi

## Enable root

``` BASH
# set root password
sudo passwd root

# enable root
sudo passwd --unlock root

# disable root
# sudo passwd --lock root

# enable root login
vi /etc/ssh/sshd_config
# update:
# permitRootLogin yes

# restart ssh
/etc/init.d/ssh restart
```

## enable colorful terminal for root

``` BASH
vi ~/.bashrc
# ~/.bashrc: executed by bash(1) for non-login shells.

# Note: PS1 and umask are already set in /etc/profile. You should not
# need this unless you want different defaults for root.
# PS1='${debian_chroot:+($debian_chroot)}\h:\w\$ '
# umask 022

# You may uncomment the following lines if you want `ls' to be colorized:
export LS_OPTIONS='--color=auto'
eval " `dircolors` "
alias ls='ls $LS_OPTIONS'
alias ll='ls $LS_OPTIONS -l'
alias l='ls $LS_OPTIONS -lA'
#
# Some more alias to avoid making mistakes:
# alias rm='rm -i'
# alias cp='cp -i'
# alias mv='mv -i'

# then exit and reconnect
```

## upgrade vi.tiny to vi

``` BASH
apt-get update
apt-get remove vim-common
apt-get install vim

# set unvisible mouse in vi
:set mouse=
```

## Update hostname of raspberry-pi

``` BASH
# update hostname
echo "edge-zero01" > /etc/hostname
sed -i $'s/127.0.1.1\t\t.*/127.0.1.1\t\tedge-zero01/' /etc/hosts
# reboot
reboot now
```

## get booting info

``` BASH
dmesg
```

## install docker

``` BASH
# install
curl -sSL https://get.docker.com | sh
# for test
docker run -d -p 8080:80 arm32v6/nginx:alpine-perl
```

## shutdown

``` BASH
shutdown -h now
```
