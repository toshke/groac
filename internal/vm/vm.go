package vm

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh"
	"gopkg.in/alessio/shellescape.v1"
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

type SSHConnnection struct {
	params *sshConnectionParams
	client *ssh.Client
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

// connect to ssh endpoint and return connection
func (conn *SSHConnnection) Connect() error {
	auth, err := publicKeyAuth(conn.params.privateKeyPath)
	if err != nil {
		return err
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
		return fmt.Errorf("failed to dial: %s", sshError)
	}
	conn.client = sshConnection
	return nil
}

func (conn *SSHConnnection) Execute(env map[string]string, command string, stdout io.Writer, stderr io.Writer) error {
	// open new session
	session, err := conn.client.NewSession()
	check(err)
	defer session.Close()

	var commandWithEnv strings.Builder
	//set session environment
	for key, value := range env {
		// assumption here that remote servers shell accepts export syntax
		commandWithEnv.WriteString(fmt.Sprintf("export %s=%s\n", key, shellescape.Quote(value)))
	}
	commandWithEnv.WriteString(command)

	//execute command by passing standard out and err streams
	session.Stdout = stdout
	session.Stderr = stderr
	return session.Run(commandWithEnv.String())

}

func (machine *Vm) Execute(cmd string) {
	fmt.Printf("Execute %v on %v", cmd, machine.Name)
	machine.ExecutedCmdsCounter++
}
