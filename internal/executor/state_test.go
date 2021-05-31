package executor

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"

	"github.com/go-playground/assert"
	"github.com/toshke/groac/internal/vm"
)

func TestSaveState(t *testing.T) {
	t.Run("verify save state", func(t *testing.T) {
		state := NewExecutorState(0)
		state.MachinesList = []vm.Vm{
			{
				Name:                "machine-1",
				ExecutedCmdsCounter: 0,
			},
			{
				Name:                "machine-2",
				ExecutedCmdsCounter: 0,
			},
		}
		state.Save()

		loadedState := NewExecutorState(0)
		loadedState.FsLoad()
		assert.Equal(t, len(loadedState.MachinesList), len(state.MachinesList))
		assert.Equal(t, loadedState.MachinesList[0].Name, state.MachinesList[0].Name)
		assert.Equal(t, loadedState.MachinesList[1].Name, state.MachinesList[1].Name)
	})
}

func TestConcurrentSaveState(t *testing.T) {
	t.Run("verify save state concurrent", func(t *testing.T) {
		// clear the state file first
		state := NewExecutorState(0)
		state.Save()
		iterationsNum := 100
		var waitGroup sync.WaitGroup
		waitGroup.Add(iterationsNum)
		c := make(chan int, iterationsNum)
		for i := 0; i < iterationsNum; i++ {
			c <- i
			go func() {
				local_counter := <-c
				fmt.Printf("Process %v\n", local_counter)
				state := NewExecutorState(rand.Intn(1000))
				state.LockDataFile()
				state.FsLoad()
				var newMachine vm.Vm
				newMachine.ExecutedCmdsCounter = 0
				newMachine.Name = fmt.Sprintf("machine-%v", local_counter)
				state.MachinesList = append(state.MachinesList, newMachine)
				state.Save()
				state.UnlockDataFile()
				waitGroup.Done()
			}()
		}

		waitGroup.Wait()
		state.FsLoad()
		assert.Equal(t, len(state.MachinesList), iterationsNum)
	})
}
