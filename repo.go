package main

import (
	"fmt"
	//"strings"
	//"bytes"
	//"os/exec"
	//"strconv"
	"errors"
)


// Initiate the repository data
func init() {
	//intialize free VNC port pool
	for i := VNC_MIN_PORT; i <= VNC_MAX_PORT; i++ {
		freeVNCPortPool = append(freeVNCPortPool,i)
	}
	//intialize free SSH port pool
	for i := SSH_MIN_PORT; i <= SSH_MAX_PORT; i++ {
		freeSSHPortPool = append(freeSSHPortPool,i)
	}
	//intialize free emulator port pool
	for i := EMULATOR_MIN_PORT; i <= EMULATOR_MAX_PORT; i = i+2 {
		freeEmulatorPortPool = append(freeEmulatorPortPool,i)
	}

	//intialize devices list
	RepoCreateDeviceInventory(Device{IMEI:"353188020902632", ADBName:"G1002de5a082",  ConnectedHostname: THIS_HOST_NAME})
	RepoCreateDeviceInventory(Device{IMEI:"353188020902633", ADBName:"G1002de5a083",  ConnectedHostname: THIS_HOST_NAME})
	RepoCreateDeviceInventory(Device{IMEI:"353188020902634", ADBName:"G1002de5a084",  ConnectedHostname: THIS_HOST_NAME})
	for i,_ := range inventories {
		if sshPort, err := RepoAllocateSSHPort(); err == nil {
			inventories[i].SSHPort = sshPort
		} else {
			panic(err)
		}
		if vncPort, err := RepoAllocateVNCPort(); err == nil {
			inventories[i].VNCPort = vncPort
		} else {
			panic(err)
		}

	}

	//fmt.Println(inventories)
	//fmt.Println(freeVNCPortPool)
	//fmt.Println(freeSSHPortPool)
	//fmt.Println(freeEmulatorPortPool)
}

const THIS_HOST_NAME = "host101"
const EMULATOR_MIN_PORT = 5554
const EMULATOR_MAX_PORT = 5584
const EMULATOR_MAX_NUM = (EMULATOR_MAX_PORT-EMULATOR_MIN_PORT)/2
const VNC_MIN_PORT = 5901
const VNC_MAX_PORT = 5920

const SSH_MIN_PORT = 5921
const SSH_MAX_PORT = 5940

const INSTALL_SCRIPT_PATH  = "/Users/Scott/master/src/github.com/tianhongbo/node/"
const INSTALL_DATA_PATH  = "/Users/Scott/master/src/github.com/tianhongbo/node/"
//VNC available port pool#
var freeVNCPortPool FreeVNCPortPool

func RepoAllocateVNCPort() (int, error) {
	for i,port := range freeVNCPortPool {
		freeVNCPortPool = append(freeVNCPortPool[:i], freeVNCPortPool[i+1:]...)
		return port, nil
	}
	return 0, errors.New("can't find the available emulator port resource.")
}

func RepoFreeVNCPort(port int) {
	freeVNCPortPool = append(freeVNCPortPool, port)
	return
}

//SSH available port pool#
var freeSSHPortPool FreeSSHPortPool

func RepoAllocateSSHPort() (int, error) {
	for i,port := range freeSSHPortPool {
		freeSSHPortPool = append(freeSSHPortPool[:i], freeSSHPortPool[i+1:]...)
		return port, nil
	}
	return 0, errors.New("can't find the available emulator port resource.")
}

func RepoFreeSSHPort(port int) {
	freeSSHPortPool = append(freeSSHPortPool, port)
	return
}


//Emulator available port pool #
var freeEmulatorPortPool FreeEmulatorPortPool

func RepoAllocateEmulatorPort() (int, error) {
	for i,port := range freeEmulatorPortPool {
		freeEmulatorPortPool = append(freeEmulatorPortPool[:i], freeEmulatorPortPool[i+1:]...)
		return port, nil
	}
	return 0, errors.New("can't find the available emulator port resource.")
}

func RepoFreeEmulatorPort(port int) {
	freeEmulatorPortPool = append(freeEmulatorPortPool, port)
	return
}

//Emulators
var emulators Emulators



func RepoFindEmulator(id int) (Emulator, error) {
	for _, e := range emulators {
		if e.Id == id {
			return e, nil
		}
	}

	return Emulator{}, errors.New("can't find the emulator.")

}

func RepoCreateEmulator(e Emulator) Emulator {
	emulators = append(emulators, e)
	return e
}

func RepoDestroyEmulator(id int) error {
	for i, e := range emulators {
		if e.Id == id {
			emulators = append(emulators[:i], emulators[i+1:]...)
			return nil
		}
	}
	return errors.New("can't find the emulator.")
}

// Mobile hubs
var hubs Hubs


func RepoFindHub(id int) (Hub,error) {
	for _, t := range hubs {
		if t.Id == id {
			return t, nil
		}
	}
	// return empty Hub if not found
	return Hub{}, errors.New("can't find the hub.")
}

func RepoCreateHub(t Hub) Hub {
	hubs = append(hubs, t)
	return t
}

func RepoDestroyHub(id int) error {
	for i, t := range hubs {
		if t.Id == id {
			hubs = append(hubs[:i], hubs[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Could not find Hub with id of %d to delete", id)
}

func RepoAttachHub(id int, connection Connection) (int, error) {
	for i, t := range hubs {
		if t.Id == id {
			for j, c := range hubs[i].Connections {
				if c.ResourceId == 0 && c.ResourceType == "" {
					hubs[i].Connections[j].ResourceId = connection.ResourceId
					hubs[i].Connections[j].ResourceType = connection.ResourceType
					return c.Port, nil
				}
			}
			return 0, fmt.Errorf("Could not find free port to attach in the hub id of %d", id)
		}

	}

	return 0, fmt.Errorf("Could not find Hub with id of %d to attach", id)
}

func RepoDetachHub(id int, connection Connection) error {
	for i, t := range hubs {
		if t.Id == id {
			for j, c := range hubs[i].Connections {
				if c.ResourceId == connection.ResourceId && c.ResourceType == connection.ResourceType {
					hubs[i].Connections[j].ResourceId = 0
					hubs[i].Connections[j].ResourceType = ""
					return nil
				}
			}
			return errors.New("can't find the resource id.")
		}

	}

	return errors.New("can't find the hub.")
}

/*
The following are devices
 */

var devices Devices

func RepoFindDevice(imei string) (Device,error) {
	for _, d := range devices {
		if d.IMEI == imei {
			return d, nil
		}
	}
	// return empty Hub if not found
	return Device{}, errors.New("can't find the device.")
}

func RepoCreateDevice(d Device) Device {
	devices = append(devices, d)
	return d
}

func RepoDestroyDevice(imei string) error {
	for i, d := range devices {
		if d.IMEI == imei {
			devices = append(devices[:i], devices[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Could not find device with id of %d to destroy", imei)
}

/*
The following are device inventories
 */

var inventories Inventroies

func RepoFindDeviceInventory(imei string) (Device,error) {
	for _, d := range inventories {
		if d.IMEI == imei {
			return d, nil
		}
	}
	// return empty Hub if not found
	return Device{}, errors.New("can't find the device in the inventory.")
}

func RepoCreateDeviceInventory(d Device) Device {
	inventories = append(inventories, d)
	return d
}

func RepoDestroyDeviceInventory(imei string) error {
	for i, d := range inventories {
		if d.IMEI == imei {
			inventories = append(inventories[:i], inventories[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Could not find device in the inventory with id of %d to destroy", imei)
}

func RepoAllocateDeviceInventory(imei string) (Device, error) {
	for _, d := range inventories {
		if d.IMEI == imei {
			fmt.Println("Device of IMEI is allocated.\n")
			return d, nil
		}
	}
	// return empty Hub if not found
	return Device{}, errors.New("can't find the device in the inventory.")
}

func RepoFreeDeviceInventory(imei string) error {
	for _, d := range inventories {
		if d.IMEI == imei {
			fmt.Println("Device of IMEI is free.")
			return nil

		}
	}
	// return empty Hub if not found
	return errors.New("can't find the device in the inventory.")
}