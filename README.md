# Radio Simulator

Created by WeiTing

## Requirement

gtp5g module

```bash
git clone https://github.com/PrinzOwO/gtp5g.git
cd gtp5g
make -j8
sudo make install
```

Can build on RPI 4 4GB with Ubuntu Server & Linux kernel 5.3.0-1017-raspi2

Cannot build on 5.3.0-1023-raspi2


## Build

```bash
# install go1.14.4
wget https://golang.org/dl/go1.14.4.linux-arm64.tar.gz
sudo tar -C /usr/local -zxvf ./go1.14.4.linux-arm64.tar.gz
echo 'export GOPATH=$HOME/go' >> ~/.bashrc
echo 'export GOROOT=/usr/local/go' >> ~/.bashrc
echo 'export PATH=$PATH:$GOPATH/bin:$GOROOT/bin' >> ~/.bashrc
source ~/.bashrc

# build go src
./build.sh

# build gtp tools
git submodule update --init
sudo apt -y install gcc cmake libmnl-dev autoconf libtool pkg-config
cd lib/upf/lib/libgtp5gnl
autoreconf -iv
./configure
make -j8
```

## Run

Run RAN

```bash
sudo ./run.sh
```

In another tty (terminal)

UE registrator

```bash
nc -v localhost 9999

imsi-2089300000003
reg # registrator
sess 1 add # add session

sess 1 del # delete session
dereg # de-registration
```

