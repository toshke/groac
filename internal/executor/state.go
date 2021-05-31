package executor

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"

	. "github.com/toshke/groac/internal/vm"
	"golang.org/x/sys/unix"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type executorState struct {
	MachinesList      []Vm `json:"machines"`
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

func withExclusiveBlocking(pathname string, callback func(*os.File)) error {
	fh, err := os.OpenFile(pathname, os.O_RDWR|os.O_CREATE, os.ModePerm)
	check(err)
	defer fh.Close()
	if err = unix.Flock(int(fh.Fd()), unix.LOCK_EX); err != nil {
		return err
	}
	callback(fh)
	return nil
}

func withSharedBlocking(pathname string, callback func(*os.File)) error {
	fh, err := os.Open(pathname)
	if err != nil {
		return err
	}
	defer fh.Close()
	if err = unix.Flock(int(fh.Fd()), unix.LOCK_SH); err != nil {
		return err
	}
	callback(fh)
	return nil
}
