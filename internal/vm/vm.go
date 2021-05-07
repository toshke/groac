package vm

import "fmt"

type Vm struct {
	Name                string
	ExecutedCmdsCounter int
}

func (machine *Vm) Execute(cmd string) {
	fmt.Printf("Execute %v on %v", cmd, machine.Name)
	machine.ExecutedCmdsCounter++
}
