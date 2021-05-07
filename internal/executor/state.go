package executor

import (
	"encoding/json"
	iofs "io/fs"
	"os"
	"path"
	"path/filepath"

	"github.com/spf13/afero"
	. "github.com/toshke/groac/internal/vm"
)

var fs = afero.NewOsFs()

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type ExecutorState struct {
	MachinesList []Vm `json:"machines"`
}

func (state *ExecutorState) json() []byte {
	jsonBytes, _ := json.Marshal(state)
	return jsonBytes
}

func (state *ExecutorState) stateFilePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		groac_path, _ := os.Executable()
		home = filepath.Dir(groac_path)
	}
	state_file := path.Join(home, ".groac", "state.json")
	return state_file
}

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		groac_path, _ := os.Executable()
		home = filepath.Dir(groac_path)
	}
	err = os.MkdirAll(path.Join(home, ".groac"), 0755)
	check(err)
}

func (state *ExecutorState) Save() {
	err := afero.WriteFile(fs, state.stateFilePath(), state.json(), iofs.ModeAppend)
	check(err)
}

func (state *ExecutorState) FsLoad() {
	bytes, err := afero.ReadFile(fs, state.stateFilePath())
	check(err)
	json.Unmarshal(bytes, state)
}
