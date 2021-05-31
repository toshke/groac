package executor

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/howeyc/fsnotify"
	"github.com/toshke/groac/internal/vm"
	"golang.org/x/sys/unix"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type executorState struct {
	MachinesList      []vm.Vm `json:"machines"`
	dataFilePath      string
	stateWriteDeferMs int
	lockedFd          *os.File
}

func NewExecutorState(writeDelayMs int) executorState {
	state := executorState{}
	state.stateWriteDeferMs = writeDelayMs
	state.initDataFile()
	return state
}

func (state *executorState) json() []byte {
	jsonBytes, _ := json.MarshalIndent(state, "", "\t")
	return jsonBytes
}

func (state *executorState) initDataFile() {
	home, err := os.UserHomeDir()
	if err != nil {
		groac_path, _ := os.Executable()
		home = filepath.Dir(groac_path)
	}
	err = os.MkdirAll(path.Join(home, ".groac"), 0755)
	check(err)

	state.dataFilePath = path.Join(home, ".groac", "state.json")
	if _, err = os.Stat(state.dataFilePath); err != nil {
		os.Create(state.dataFilePath)
	}
}

// Allowing both for explicit and implicit file locks. Implict will take place
// explicit lock hasn't been requested.

func (state *executorState) save() {
	state.lockedFd.Truncate(0)
	state.lockedFd.Seek(0, 0)
	_, err := state.lockedFd.Write(state.json())
	time.Sleep(time.Duration(state.stateWriteDeferMs) * time.Millisecond)
	check(err)
}

func (state *executorState) Save() {
	if state.lockedFd != nil {
		state.save()
	} else {
		state.LockDataFile()
		state.save()
		state.UnlockDataFile()
	}
}

func (state *executorState) fsLoad() {
	reader := io.Reader(state.lockedFd)
	bytes, err := ioutil.ReadAll(reader)
	check(err)
	json.Unmarshal(bytes, state)
}

func (state *executorState) FsLoad() {
	if state.lockedFd != nil {
		state.fsLoad()
	} else {
		state.LockDataFile()
		state.fsLoad()
		state.UnlockDataFile()
	}
}

func (state *executorState) LockDataFile() {
	fh, err := os.OpenFile(state.dataFilePath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	check(err)
	err = unix.Flock(int(fh.Fd()), unix.LOCK_EX)
	check(err)
	state.lockedFd = fh
}

func (state *executorState) UnlockDataFile() {
	if state.lockedFd != nil {
		err := state.lockedFd.Close()
		if err != nil {
			check(err)
		}
		state.lockedFd = nil
	}
}

func (state *executorState) EnableReload() {
	watcher, err := fsnotify.NewWatcher()
	check(err)

	// done := make(chan bool)

	// Process events
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				log.Println("event:", ev)
				state.FsLoad()
			case err := <-watcher.Error:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Watch(state.dataFilePath)
	check(err)

	// Hang so program doesn't exit
	// <-done

	/* ... do stuff ... */
	// watcher.Close()

}
