package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"os/exec"
	"strings"
	"bytes"
	"time"
	"github.com/gorilla/mux"

)



func Index(w http.ResponseWriter, r *http.Request) {

	cmdStr := "emulator -avd  Android_2.2 -no-window -verbose -no-boot-anim -noskin"
	//cmdStr := "android list target"
	parts := strings.Fields(cmdStr)
	command := parts[0]
	args := parts[1:len(parts)]
	fmt.Println(command, args)

	//parts := strings.Fields(cmd)
	//head := parts[0]
	//parts = parts[1:len(parts)]
	//fmt.Println(head, parts)

	cmd := exec.Command(command, args...)
	randomBytes := &bytes.Buffer{}
	cmd.Stdout = randomBytes

	// Start command asynchronously
	err := cmd.Start()
	fmt.Println("A emulator is launched: %s \n %s", err, randomBytes.String())
	fmt.Fprintf(w, "Successful!\n")
}

// Emulator
// /emulators GET
func EmulatorIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(emulators); err != nil {
		panic(err)
	}
}

// /emulators POST
func EmulatorCreate(w http.ResponseWriter, r *http.Request) {
	var port int
	var e Emulator

	//emulators.showAll()
	
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	
	fmt.Println("the body of request is: ", r.Header)
	fmt.Println("the body of request is: ", string(body))

	if err := json.Unmarshal(body, &e); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}

	fmt.Println("the Json from request is: ", e)

	if port, err = RepoAllocateEmulatorPort(); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}

	e.init(port)
	e.start()

	RepoCreateEmulator(e)

	fmt.Println("Create emulator successfully: ", e.Cmd)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(e); err != nil {
		panic(err)
	}

}


// /emulators/{id} DELETE
func EmulatorDestroy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var id int
	var err error

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		panic(err)
	}

	if e,err := RepoFindEmulator(id); err == nil {

		fmt.Println("Delete emulator successfully: ", e.Cmd)
		e.stop()
		RepoFreeEmulatorPort(e.Port)
		RepoDestroyEmulator(id)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	return
}

// /emulators/{id} GET
func EmulatorShow(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	var id int
	var err error
	var e Emulator

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		panic(err)
	}

	if e,err = RepoFindEmulator(id); err == nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(e); err != nil {
			panic(err)
		}
	}
	return
}


func HubIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(hubs); err != nil {
		panic(err)
	}
}

func HubShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var hubId int
	var err error
	if hubId, err = strconv.Atoi(vars["id"]); err != nil {
		panic(err)
	}
	hub,err := RepoFindHub(hubId)

	if err == nil {
		//Find the hub
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(hub); err != nil {
			panic(err)
		}
	} else {
		// If we didn't find a hub, 404
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
			panic(err)
		}
	}


	return

}

func HubDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var hubId int
	var err error
	if hubId, err = strconv.Atoi(vars["id"]); err != nil {
		panic(err)
	}
	err = RepoDestroyHub(hubId)
	if err == nil {
		// Find the hub and deleted successfully
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

	} else {
		// If we didn't find it, 404
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
			panic(err)
		}

	}

	return
}


/*
Test with this curl command:

curl -H "Content-Type: application/json" -d '{"name":"New Hub"}' http://localhost:8080/hubs

*/
func HubCreate(w http.ResponseWriter, r *http.Request) {
	var hub Hub
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &hub); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	//Initialize the connections

	for i := 0; i < hub.PortNum; i++ {
		hub.Connections = append(hub.Connections, Connection{i, "", 0})
	}
	//Initialize the start time
	hub.StartTime = time.Now()
	

	t := RepoCreateHub(hub)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(t); err != nil {
		panic(err)
	}
}

func HubAttach(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var hubId int
	var connection Connection
	var err error

	if hubId, err = strconv.Atoi(vars["id"]); err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &connection); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}

	if i ,err := RepoAttachHub(hubId, connection); err == nil {
		// attach successfully

		// manipulate the device network
		fmt.Println("Turn on the intenet connection ")

		// respond
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		t := AttachRsp{connection.ResourceType, connection.ResourceId, hubId, i}
		if err := json.NewEncoder(w).Encode(t); err != nil {
			panic(err)
		}

	} else {
		// attach failed
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(jsonErr{Code: 422, Text: "No available port."}); err != nil {
			panic(err)
		}

	}

	return

}

func HubDetach(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var hubId int
	var err error

	connection := Connection{}

	if hubId, err = strconv.Atoi(vars["id"]); err != nil {
		panic(err)
	}


	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &connection); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}

	if err := RepoDetachHub(hubId, connection); err == nil {
		// Detach successfully

		// manipulate the device network
		fmt.Println("Turn off the intenet connection ")

		// respond
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

	} else {
		// detach failed
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
			panic(err)
		}

	}

	return
}

/*
Device Operation
 */

// /devices GET
func DeviceIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(devices); err != nil {
		panic(err)
	}
}

// /device POST
func DeviceCreate(w http.ResponseWriter, r *http.Request) {
	var d,inventory Device

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &d); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}

	if  inventory, err = RepoAllocateDeviceInventory(d.IMEI); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}

	d.ADBName = inventory.ADBName
	d.ConnectedHostname = inventory.ConnectedHostname
	d.SSHPort = inventory.SSHPort
	d.VNCPort = inventory.VNCPort

	//Initialize the start time
	d.StartTime = time.Now()

	d = RepoCreateDevice(d)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(d); err != nil {
		panic(err)
	}

	return
}

// /devices/{imei} GET
func DeviceShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var imei string
	var err error
	if imei = vars["imei"]; imei == "" {
		panic(err)
	}
	if d,err := RepoFindDevice(imei); err == nil {
		//Find the device
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(d); err != nil {
			panic(err)
		}
	} else {
		// If we didn't find a device, 404
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
			panic(err)
		}
	}


	return

}

// /devices/{imei}
func DeviceDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var imei string
	var err error
	if imei = vars["imei"]; imei == "" {
		panic(err)
	}
	if err = RepoDestroyDevice(imei); err == nil {
		// Find the hub and deleted successfully
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

	} else {
		// If we didn't find it, 404
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
			panic(err)
		}

	}

	return
}

