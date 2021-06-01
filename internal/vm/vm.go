package vm

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"golang.org/x/crypto/ssh"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Vm struct {
	Name                string `json:"name"`
	ExecutedCmdsCounter int    `json:"counter"`
}

type VmConnection interface {
	Connect()
}

type sshConnectionParams struct {
	username       string
	privateKeyPath string
	publicKeyPath  string
	hostname       string
	port           int16
}

func newSshConnectionParams() *sshConnectionParams {
	var params sshConnectionParams
	params.username = "gitlab-runner"
	params.port = 22
	params.hostname = "127.0.0.1"
	return &params
}

func (params *sshConnectionParams) initPrivateKey() {
	home, err := os.UserHomeDir()
	if len(params.privateKeyPath) == 0 {
		if err != nil {
			groac_path, _ := os.Executable()
			home = filepath.Dir(groac_path)
		}
		err = os.MkdirAll(path.Join(home, ".groac"), 0755)
		check(err)
		params.privateKeyPath = path.Join(home, ".groac", "vms_key.pem")
	}
	params.publicKeyPath = fmt.Sprintf("%s.pub", params.privateKeyPath)
	// TODO: only if not exists
	if _, err = os.Stat(params.privateKeyPath); err != nil {
		check(generateKeyPair(params.privateKeyPath, params.publicKeyPath))
	}
}

type SSHConnnection struct {
	params *sshConnectionParams
}

// connect to ssh endpoint and return connection
func (conn *SSHConnnection) Connect() (*ssh.Client, error) {
	auth, err := publicKeyAuth(conn.params.privateKeyPath)
	if err != nil {
		return nil, err
	}
	sshConfig := &ssh.ClientConfig{
		User: conn.params.username,
		Auth: []ssh.AuthMethod{
			auth,
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	destination := fmt.Sprintf("%s:%d", conn.params.hostname, conn.params.port)
	sshConnection, sshError := ssh.Dial("tcp", destination, sshConfig)
	if sshError != nil {
		return nil, fmt.Errorf("failed to dial: %s", sshError)
	}
	return sshConnection, nil
}

func (machine *Vm) Execute(cmd string) {
	fmt.Printf("Execute %v on %v", cmd, machine.Name)
	machine.ExecutedCmdsCounter++
}
