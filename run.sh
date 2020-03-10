#!/bin/bash
exe() { echo "\$ $@" ; "$@" ; }

TUN="uetun"
LISTEN_ADDR="127.0.0.2"
UE_SUBNET="60.60.0.1/16"
SERVER_SUBNET="60.60.0.1/16"


# Setup tunnel interface 
exe sudo ip tunnel add ${TUN} mode ipip remote ${LISTEN_ADDR} local 127.0.0.1
exe sudo ip link set ${TUN} up
exe sudo ip addr add ${UE_SUBNET} peer ${SERVER_SUBNET} dev ${TUN}

sudo ./bin/simulator & 
PID=$!

trap "sudo kill -SIGKILL $PID" SIGINT SIGTERM

echo starting
# commands to start your services go here

wait

# commands to shutdown your services go here
echo exited

# Del tunnel interface
exe sudo ip link del ${TUN}