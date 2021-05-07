package remote

import (
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/toshke/groac/internal/vm"
)

func ExampleRemoteExecute() {
	vm := &vm.Vm{
		Name: "test-vm",
	}
	RemoteExecute(vm, "docker info")
	// Output: Execute docker info on test-vm
}

func TestRemoteExecute(t *testing.T) {
	t.Run("verify increment", func(t *testing.T) {
		vm := &vm.Vm{
			Name: "test-vm",
		}
		RemoteExecute(vm, "docker info")
		got := vm.ExecutedCmdsCounter
		assert.Equal(t, got, 1)
	})

}
