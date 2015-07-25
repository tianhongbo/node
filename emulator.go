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

type FreeEmulatorPorts []int

type Emulator struct {
	Id			int 		`json:"id"`
	Name		string 		`json:"name"`
	ConnectedHostname string `json: "connected_hostname"`
	VNCPort 	int 		`json: "vnc_port"`
	SSHPort 	int 		`json: "ssh_port"`
	ADBName 	string 		`json: "-"`
	Port		int 		`json:"port"`
	StartTime	time.Time	`json:"start_time"`
	StopTime	time.Time 	`json:"stop_time"`
	Cmd			*exec.Cmd 	`json:"-"`
}


type Emulators []Emulator


type ApiStartEmulatorResponse struct {
	Id	uint64	`json:"id"`
	Port	int	`json:"port"`
	
}


type ApiShowEmulatorResponse struct {
	Id		int		`json:"id"`
	Name		string		`json:"name"`
	Port		int		`json:"port"`
	StartTime	time.Time	`json:"starttime"`
	StopTime	time.Time	`json:"stoptime"`
}


func (emu *Emulator) start() {
	err := emu.Cmd.Start()
	emu.StartTime = time.Now()
	fmt.Println("A emulator is successfully launched at port: ", err, emu.Port)
	
}

func (emu *Emulator) stop() {
	emu.StopTime = time.Now()
	err := emu.Cmd.Process.Signal(os.Kill)
	err = emu.Cmd.Wait()

	fmt.Println("A emulator is successfully stopped. The result of cmd is: ", err) 
}

func (emu *Emulator) init(port int) {
	emu.StartTime = time.Now()	//startTime	time.Time
	emu.StopTime = time.Time{}	//stopTime	time.Time
                
	emu.Port = port

	//init cmd emulator64-arm -avd myandroid -no-window -verbose -no-boot-anim -noskin
	cmdStr := "emulator64-arm -avd android-api-22 -no-window -verbose -no-boot-anim -noskin -port " + strconv.Itoa(int(emu.Port))
	parts := strings.Fields(cmdStr)
	head := parts[0]
	parts = parts[1:len(parts)]
        
	//fmt.Println(head, parts)
	emu.Cmd = exec.Command(head, parts...)
	randomBytes := &bytes.Buffer{}
	emu.Cmd.Stdout = randomBytes

	fmt.Println(emu.Cmd)
	fmt.Println(emu)
}
