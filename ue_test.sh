#!/bin/bash
#
# Bash must have been compiled with this ability: --enable-net-redirections
# The device files below do not actually exist.
# Use /dev/udp for UDP sockets

# msg=$'asfasdf\nsdfasdfasdf\n'


# IFS=$'\n' read -rd '' -a array <<<"$msg"
# for i in "${array[@]}"; do
#     echo ${i}
# done

SESS_FORMAT=$'\[SESSION\] ID=([0-9]+),DNN=([^,]+),SST=([0-9]+),SD=([0-9]+),UEIP=([^,]+),ULAddr=([^,]+),ULTEID=([0-9]+),DLAddr=([^,]+),DLTEID=([0-9]+)'
# if [[ $test_string =~ $SESS_FORMAT ]]; then echo "DNN=${BASH_REMATCH[1]},UEIP=${BASH_REMATCH[4]}"; fi

# TUN="uetun"
# TUN_ADDR="60.60.0.1"


HOST="127.0.0.1"
PORT="9999"
SUPI="imsi-2089300007487"
ALIVE=false
TIME=0

ID=10


function check_error() {
    if [[ "$1" == *"[ERROR]"* ]] && [[ "$1" == *"FAIL"* ]]; then
        exit 1
    fi
}

send_msg() { 
    echo "\$ $1" 
    echo "$1" >&$2
}
read_msg() {
    read -r msg_in <&$1
    check_error "$msg_in"
    echo "$msg_in"
}
get_ueip(){
    if [[ $1 =~ $SESS_FORMAT ]]
    then 
        echo "${BASH_REMATCH[5]}"
    fi 
}

exec 3<>/dev/tcp/${HOST}/${PORT}

read -r msg_in <&3
echo $msg_in
# send SUPI

send_msg "$SUPI" 3
read_msg 3

# Register
send_msg "reg" 3
read_msg 3

# ADD Session
msg_out="sess 10 add"
[ -n "$ID" ] && msg_out="sess $ID add"
[ -n "$SLICE" ] && msg_out="$msg_out"" ${SLICE}"
send_msg "$msg_out" 3
msg_in=$(read_msg 3)
echo "$msg_in"

# Add Ip in tun dev
# UEIP=$(get_ueip "$msg_in")
# if [ -n "${UEIP}" ] && [ "${UEIP}" != ${TUN_ADDR} ]
# then
#     sudo ip addr add ${UEIP} dev ${TUN}
# fi


if [ $TIME -gt 0 ]
then 
    echo "Wait $TIME seconds"
    sleep $TIME &
else 
    sleep infinity & 
fi

function terminate(){
    if $ALIVE;
    then
        # send rel pdu sess
        send_msg "$(echo "${msg_out}" | sed -e "s/add/del/g")" 3
        read_msg 3
    else 
        # send del reg
        send_msg "dereg" 3
        read_msg 3
    fi
    if [ -n "${UEIP}" ] && [ "${UEIP}" != ${TUN_ADDR} ]
    then
        sudo ip addr del ${UEIP} dev ${TUN}
    fi
    exit 1
}

trap terminate SIGINT
wait 

terminate