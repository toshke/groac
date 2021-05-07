package vm

import "fmt"

type Vm struct {
	Name                string `json:"name"`
	ExecutedCmdsCounter int    `json:"counter"`
}

func (machine *Vm) Execute(cmd string) {
	fmt.Printf("Execute %v on %v", cmd, machine.Name)
	machine.ExecutedCmdsCounter++
}
