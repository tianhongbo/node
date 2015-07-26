#!/bin/bash

# ADB name, vnc_port, ssh_port
adb_name=$1
vnc_port=$2
ssh_port=$3
echo "adb_name=$adb_name, vnc_port=$vnc_port, ssh_port=$ssh_port"

#waiting for device online
adb wait-for-device

#waiting for device booting
A=$(adb -s $adb_name shell getprop sys.boot_completed | tr -d '\r')
while [ "$A" != "1" ]; do
        sleep 2
        A=$(adb -s $adb_name shell getprop sys.boot_completed | tr -d '\r')
done

#unlock emulator screen
adb -s $adb_name shell input keyevent 82

#disconnect Internet connection
adb -s $adb_name shell svc data disable
adb -s $adb_name shell svc wifi disable

#install, start vnc server service on the emulator
adb -s $adb_name push androidvncserver /data/
adb -s $adb_name shell chmod 755 /data/androidvncserver
adb -s $adb_name forward tcp:$vnc_port tcp:5901
adb -s $adb_name shell /data/androidvncserver -k "/dev/input/event0" -t "/dev/input/event0"

