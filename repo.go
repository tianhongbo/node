package main

import (
	"fmt"
	//"strings"
	//"bytes"
	//"os/exec"
	//"strconv"
)

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



var currentId int

var todos Todos


// Give us some seed data
func init() {
	RepoCreateTodo(Todo{Name: "Write presentation"})
	RepoCreateTodo(Todo{Name: "Host meetup"})
}

func RepoFindTodo(id int) Todo {
	for _, t := range todos {
		if t.Id == id {
			return t
		}
	}
	// return empty Todo if not found
	return Todo{}
}

//this is bad, I don't think it passes race condtions
func RepoCreateTodo(t Todo) Todo {
	currentId += 1
	t.Id = currentId
	todos = append(todos, t)
	return t
}

func RepoDestroyTodo(id int) error {
	for i, t := range todos {
		if t.Id == id {
			todos = append(todos[:i], todos[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Could not find Todo with id of %d to delete", id)
}
