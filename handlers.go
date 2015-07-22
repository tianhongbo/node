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

func StartEmulator(w http.ResponseWriter, r *http.Request) {
	
	var req ApiStartEmulatorRequest
	var res ApiStartEmulatorResponse
	
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

	if err := json.Unmarshal(body, &req); err != nil || req.Id == 0 {
		fmt.Println("the Json from request is: ", req)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	fmt.Println("the Json from request is: ", req)

	emulator := emulators.allocate()

	//fmt.Println("the emulator before start is: ", emulator)

	if emulator != nil {
		emulator.id = req.Id
		emulator.start()
		//fmt.Println(emulator)
		fmt.Println("Launch emulator successfully.")
		res.Id = emulator.id
		res.Port = emulator.port
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			panic(err)
		}	
	
	} else {
	
		fmt.Println("Launch failure.")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	//fmt.Println("the emulator aftre start is: ", emulator)

	emulators.showAll()
}

func StopEmulator(w http.ResponseWriter, r *http.Request) {
	
	var req ApiStopEmulatorRequest
	var res ApiStopEmulatorResponse

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	fmt.Println("the head of request is: ", r.Header)
	fmt.Println("the body of request is: ", string(body))
	
	if err := json.Unmarshal(body, &req); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	fmt.Println("the json of request is: ", req)
	
	emulator := emulators.getEmulator(req.Id)
	if emulator != nil {
		emulator.stop()
		res.Id = emulator.id
		res.StartTime = emulator.startTime
		res.StopTime = emulator.stopTime
		emulator.init(emulator.port)
		
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		
		if err := json.NewEncoder(w).Encode(res); err != nil {
			panic(err)
		}

	} else {
		fmt.Println("can't find this emulator. d%", req.Id)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	emulators.showAll()
}

func ShowEmulator(w http.ResponseWriter, r *http.Request) {
	
	var res ApiShowEmulatorResponse
	var err error

	str := r.URL.Query().Get("id")
	if str != "" { 
	// request carries "id" parameter
		emulatorId,_ := strconv.ParseUint(str, 10, 64)
		fmt.Println("the emulator id in the request is: ", emulatorId)
	
		emulator := emulators.getEmulator(emulatorId)
		if emulator != nil {
			res.Id = emulator.id
			res.Name = emulator.name
			res.Port = emulator.port
			res.StartTime = emulator.startTime
			res.StopTime = emulator.stopTime
		
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)
		
			if err := json.NewEncoder(w).Encode(res); err != nil {
				panic(err)
			}	

		} else {
			fmt.Println("can't find this emulator. d%", emulatorId)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(422) // unprocessable entity
			if err := json.NewEncoder(w).Encode(err); err != nil {
				panic(err)
			}
		}
	} else {
	// no "id" para in the request from client

		var ress []ApiShowEmulatorResponse
		for _, e := range emulators.emulators {
			if e.status == true {
				res := ApiShowEmulatorResponse {e.id, 
								e.name, 
								e.port, 
								e.startTime,
								e.stopTime,
								}
				ress = append(ress, res)
			}
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(ress); err != nil {
			panic(err)
		}

	}

	emulators.showAll()
}



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


