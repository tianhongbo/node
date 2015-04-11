package main

import (
	"fmt"
	"strings"
	"time"
	"os"
	"os/exec"
	"strconv"
	"bytes"
)

type Emulator struct {
	id		uint64
	name		string
	status		bool	//true: used, false: idle
	startTime	time.Time
	stopTime	time.Time
	cmd		*exec.Cmd
	port		uint16 
}

type ApiStartEmulatorRequest struct {
	Id	uint64	`json:"id"`
	Name	string	`json:"name"`	
}

type ApiStartEmulatorResponse struct {
	Id	uint64	`json:"id"`
	Port	uint16	`json:"port"`
	
}

type ApiStopEmulatorRequest struct {
	Id	uint64	`json:"id"`
	
}

type ApiStopEmulatorResponse struct {
	Id		uint64		`json:"id"`
	StartTime	time.Time	`json:"starttime"`
	StopTime	time.Time	`json:"stoptime"`
}

type ApiShowEmulatorRequest struct {
	//Id	uint64	`json:"id"`
	
}

type ApiShowEmulatorResponse struct {
	Id		uint64		`json:"id"`
	Name		string		`json:"name"`
	Port		uint16		`json:"port"`
	StartTime	time.Time	`json:"starttime"`
	StopTime	time.Time	`json:"stoptime"`
}


func (emu *Emulator) start() {
	err := emu.cmd.Start()
	emu.startTime = time.Now()
	fmt.Println("A emulator is successfully launched at port: ", err, emu.port) 
	
}

func (emu *Emulator) stop() {
	emu.stopTime = time.Now()
	err := emu.cmd.Process.Signal(os.Kill)
	err = emu.cmd.Wait()

	fmt.Println("A emulator is successfully stopped. The result of cmd is: ", err) 
}

func (emu *Emulator) init(port uint16) {
	emu.id = 0		//id		uint64
	emu.name = ""		//name		string
	emu.status = false	//status		bool	//true: used, false: idle
	emu.startTime = time.Time{}	//startTime	time.Time
	emu.stopTime = time.Time{}	//stopTime	time.Time
                
        emu.port = port

        //init cmd
        cmdStr := "emulator -avd Android_2.2 -no-window -verbose -no-boot-anim -noskin -port " + strconv.Itoa(int(emu.port))
        parts := strings.Fields(cmdStr)
        head := parts[0]
        parts = parts[1:len(parts)]
        
	//fmt.Println(head, parts)
        emu.cmd = exec.Command(head, parts...)
        randomBytes := &bytes.Buffer{}
        emu.cmd.Stdout = randomBytes

        fmt.Println(emu.cmd)
        fmt.Println(emu)
}
