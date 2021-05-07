package executor

import (
	"testing"

	"github.com/go-playground/assert"
	"github.com/toshke/groac/internal/vm"
)

func TestSaveState(t *testing.T) {
	t.Run("verify save state", func(t *testing.T) {
		state := &ExecutorState{
			MachinesList: []vm.Vm{
				{
					Name:                "machine-1",
					ExecutedCmdsCounter: 0,
				},
				{
					Name:                "machine-2",
					ExecutedCmdsCounter: 0,
				},
			},
		}
		state.Save()

		loadedState := &ExecutorState{}
		loadedState.FsLoad()
		assert.Equal(t, len(loadedState.MachinesList), len(state.MachinesList))
		assert.Equal(t, loadedState.MachinesList[0].Name, state.MachinesList[0].Name)
		assert.Equal(t, loadedState.MachinesList[1].Name, state.MachinesList[1].Name)
	})
}
