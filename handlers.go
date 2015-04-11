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
//	"time"
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
	
	//fmt.Println("the body of request is: ", string(body))
	
	if err := json.Unmarshal(body, &req); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	//fmt.Println("the Json from request is: ", req)

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
	//fmt.Println("the body of request is: ", string(body))
	if err := json.Unmarshal(body, &req); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

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
	
	var req ApiShowEmulatorRequest
	var res ApiShowEmulatorResponse

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &req); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	emulator := emulators.getEmulator(req.Id)
	if emulator != nil {
		res.Id = emulator.id
		res.Name = emulator.name
		res.Port = emulator.port
		res.StartTime = emulator.startTime
		res.StopTime = emulator.stopTime
		
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


func TodoIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(todos); err != nil {
		panic(err)
	}
}

func TodoShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var todoId int
	var err error
	if todoId, err = strconv.Atoi(vars["todoId"]); err != nil {
		panic(err)
	}
	todo := RepoFindTodo(todoId)
	if todo.Id > 0 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(todo); err != nil {
			panic(err)
		}
		return
	}

	// If we didn't find it, 404
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
		panic(err)
	}

}

/*
Test with this curl command:

curl -H "Content-Type: application/json" -d '{"name":"New Todo"}' http://localhost:8080/todos

*/
func TodoCreate(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &todo); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	t := RepoCreateTodo(todo)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(t); err != nil {
		panic(err)
	}
}
