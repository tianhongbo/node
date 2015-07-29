package main

import (
	"fmt"
	"strings"
	"time"
	"os/exec"
	"strconv"

	"bytes"
)

type FreeEmulatorPortPool []int

type Emulator struct {
	Id			int 		`json:"id"`
	Name		string 		`json:"name"`
	ConnectedHostname string `json: "connected_hostname"`
	VNCPort 	int 		`json: "vnc_port"`
	SSHPort 	int 		`json: "ssh_port"`
	ADBName 	string 		`json: "-"`
	EmulatorPort int 		`json:"port"`
	StartTime	time.Time	`json:"start_time"`
	StopTime	time.Time 	`json:"stop_time"`
	Cmd			*exec.Cmd 	`json:"-"`
	CmdInit		*exec.Cmd 	`json:"-"`	//Initializing Script exectued after emulator start
	CmdAttach		*exec.Cmd 	`json:"-"`	//Enable network connection Script
	CmdDetach		*exec.Cmd 	`json:"-"`	//Disable network connection script
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
	emu.Cmd.Start()
	emu.CmdInit.Start()

	fmt.Println("A emulator is successfully launched at port: ", emu.EmulatorPort)
	
}

func (emu *Emulator) stop() {
	var err error

	emu.StopTime = time.Now()

	//err = emu.Cmd.Process.Signal(os.Kill)
	//err = emu.Cmd.Wait()
	//fmt.Println(emu.Cmd.Stdout)

	err = emu.Cmd.Process.Kill()
	emu.Cmd.Wait()


	emu.CmdInit.Process.Kill()

	err = emu.CmdInit.Wait()

//	fmt.Println(emu.CmdInit.Stdout)
	//err = syscall.Kill(emu.CmdInit.Process.Pid, 9)


	fmt.Println("A emulator is successfully stopped.", err)
}

func (emu *Emulator) initCmd() {


	//init cmd emulator64-arm -avd myandroid -no-window -verbose -no-boot-anim -noskin
	cmdStr := "emulator64-arm -avd android-api-10-"+strconv.Itoa(int(emu.EmulatorPort))+" -wipe-data -no-window -no-boot-anim -noskin -port " + strconv.Itoa(int(emu.EmulatorPort))
	parts := strings.Fields(cmdStr)
	head := parts[0]
	parts = parts[1:len(parts)]
	emu.Cmd = exec.Command(head, parts...)
	randomBytes := &bytes.Buffer{}
	emu.Cmd.Stdout = randomBytes

	//init cmd install.sh emulator-5566 vnc_port ssh_port
	cmdStr2 := INSTALL_SCRIPT_PATH + "install.sh " + emu.ADBName + " " + strconv.Itoa(emu.VNCPort) + " " + strconv.Itoa(emu.SSHPort) + " " + strconv.Itoa(emu.EmulatorPort)

	parts2 := strings.Fields(cmdStr2)
	head2 := parts2[0]
	parts2 = parts2[1:len(parts2)]
	emu.CmdInit = exec.Command(head2, parts2...)

	randomBytes2 := &bytes.Buffer{}
	emu.CmdInit.Stdout = randomBytes2

}

func (emu *Emulator) attach() {

	cmdStr := INSTALL_SCRIPT_PATH + "attach.sh " + emu.ADBName
	parts := strings.Fields(cmdStr)
	head := parts[0]
	parts = parts[1:len(parts)]
	cmd := exec.Command(head, parts...)
	//randomBytes := &bytes.Buffer{}
	//cmd.Stdout = randomBytes
	err := cmd.Run()
	fmt.Println("The emualator was attached, error code = ", err)
}

func (emu *Emulator) detach() {

	cmdStr := INSTALL_SCRIPT_PATH + "detach.sh " + emu.ADBName
	parts := strings.Fields(cmdStr)
	head := parts[0]
	parts = parts[1:len(parts)]
	cmd := exec.Command(head, parts...)
	//randomBytes := &bytes.Buffer{}
	//cmd.Stdout = randomBytes
	err := cmd.Run()
	fmt.Println("The emualator was detached, error code = ", err)
}