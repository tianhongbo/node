#!/bin/bash
#exit when go main() send ok.Interrupt signal to this process
trap 'echo "Exit 2(os.Interrupt signal detected... vncserver_id=$vncserver_pid, novnc_pid=$novnc_pid"; kill -9 $vncserver_pid; kill -9 $novnc_pid; exit 0' 2

#Please make sure two paths configuration
vnc_server_path=/Users/Scott/master/src/github.com/tianhongbo/node
novnc_path=/Users/Scott/noVNC

# ADB name, vnc_port, ssh_port
adb_name=$1
vnc_port=$2
ssh_port=$3
emulator_port=$4

echo "adb_name=$adb_name, vnc_port=$vnc_port, ssh_port=$ssh_port, emulator_port=$emulator_port"

#Create AVD
echo no | android -s create avd -n android-api-10-$emulator_port -t android-10 --abi default/armeabi

#Start Emulator
#emulator64-arm -avd android-api-10-$emulator_port -wipe-data -no-window -no-boot-anim -noskin -port $emulator_port&

#waiting for device online
adb -s $adb_name wait-for-device

#waiting for device booting
A=$(adb -s $adb_name shell getprop sys.boot_completed | tr -d '\r')
while [ "$A" != "1" ]; do
        sleep 2
        A=$(adb -s $adb_name shell getprop sys.boot_completed | tr -d '\r')
done

#unlock emulator screen
#adb -s $adb_name shell input keyevent 82

#disconnect Internet connection
#adb -s $adb_name shell svc data disable
#adb -s $adb_name shell svc wifi disable
adb -s $adb_name shell setprop net.dns1 0.0.0.0

#configure for SSH
adb -s $adb_name forward tcp:$ssh_port tcp:22

#install, start vnc server service on the emulator
adb -s $adb_name push $vnc_server_path/androidvncserver /data/
adb -s $adb_name shell chmod 755 /data/androidvncserver
adb -s $adb_name forward tcp:$vnc_port tcp:5901
adb -s $adb_name shell /data/androidvncserver -k "/dev/input/event0" -t "/dev/input/event0"&
vncserver_pid=$!

cd $novnc_path
$novnc_path/utils/launch.sh --listen $vnc_port --vnc localhost:$vnc_port --web $novnc_path&
novnc_pid=$!

while [ 1 ]
do
        sleep 2
done