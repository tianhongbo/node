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
	//RepoCreateHub(Hub{Name: "SJSU_GUEST_WIFI"})
	//RepoCreateHub(Hub{Name: "AT&T_LTE"})
	RepoCreateDeviceInventory(Device{IMEI:"353188020902632", ADBName:"G1002de5a082", VNCPort: 5900, SSHPort: 22, ConnectedHostname: THIS_HOST_NAME})
	RepoCreateDeviceInventory(Device{IMEI:"353188020902633", ADBName:"G1002de5a083", VNCPort: 5901, SSHPort: 23, ConnectedHostname: THIS_HOST_NAME})
	RepoCreateDeviceInventory(Device{IMEI:"353188020902634", ADBName:"G1002de5a084", VNCPort: 5902, SSHPort: 24, ConnectedHostname: THIS_HOST_NAME})
}

const THIS_HOST_NAME = "host101"
//Emulators
const EMULATOR_MIN_PORT = 5554
const EMULATOR_MAX_PORT = 5584
const EMULATOR_NUM = (EMULATOR_MAX_PORT-EMULATOR_MIN_PORT)/2

type Emulators struct {
	emulators [EMULATOR_NUM]Emulator
}

var emulators Emulators

func (emus *Emulators) init() {
	for i := 0; i < EMULATOR_NUM; i++ {
		emus.emulators[i].init(uint16(EMULATOR_MIN_PORT + i*2)) 
	}
}

func (emus *Emulators) allocate()  (*Emulator) {
	for i := 0; i < EMULATOR_NUM; i++ {
		if emus.emulators[i].status == false {
			emus.emulators[i].status = true //set used flag "true"
			//fmt.Println(i)
			//fmt.Println(emus.emulators[i])
			return &emus.emulators[i]
		}
	}
	return nil 
}

func (emus *Emulators) showAll() {
	for i, emu := range emus.emulators {
		fmt.Println("the i = %d emulator: ", i, emu)
	}
}


func (emus *Emulators) getEmulator(id uint64)  (*Emulator) {
	if id == 0 {return nil}
	for i,_ := range emus.emulators {
		if emus.emulators[i].id == id {
			fmt.Println(emus.emulators[i])
			return &emus.emulators[i]
		}
	}
	return nil
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