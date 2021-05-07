package remote

import (
	"github.com/toshke/groac/internal/vm"
)

func RemoteExecute(machine *vm.Vm, cmd string) {
	machine.Execute(cmd)
}
