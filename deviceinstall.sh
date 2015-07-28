#!/bin/bash

vnc_server_path=/Users/Scott/master/src/github.com/tianhongbo/node
novnc_path=/Users/Scott/noVNC

# ADB name, device_ip, vnc_port, ssh_port
adb_name=$1
device_ip=$2
vnc_port=$3
ssh_port=$4
echo "adb_name=$adb_name, device_ip=$device_ip, vnc_port=$vnc_port, ssh_port=$ssh_port"

#waiting for device online
adb -s $adb_name wait-for-device

#waiting for device booting
A=$(adb -s $adb_name shell getprop sys.boot_completed | tr -d '\r')
while [ "$A" != "1" ]; do
        sleep 1
        A=$(adb -s $adb_name shell getprop sys.boot_completed | tr -d '\r')
done


#disconnect Internet connection
#adb -s $adb_name shell 'su -c "svc wifi disable"'
#adb -s $adb_name shell 'su -c "svc data disable"'

#adb -s $adb_name shell setprop net.dns1 0.0.0.0
adb -s $adb_name shell 'su -c "setprop net.dns1 0.0.0.0"'

#start vnc proxy on the host
#/Users/Scott/noVNC/utils/launch.sh --listen 5910 --vnc 192.168.1.16:5901 --web /Users/Scott/noVNC
cd $novnc_path
$novnc_path/utils/launch.sh --listen $vnc_port --vnc $device_ip:5901 --web $novnc_path&


