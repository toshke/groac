package vm

import (
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/go-playground/assert/v2"
)

// test private key generation
func TestNewSSHParams(t *testing.T) {
	t.Run("verify save state", func(t *testing.T) {
		sshConnectionParams := newSshConnectionParams()
		sshConnectionParams.initPrivateKey()

		auth, err := publicKeyAuth(sshConnectionParams.privateKeyPath)
		assert.NotEqual(t, auth, nil)
		assert.Equal(t, err, nil)
	})
}

func TestSSHConnect(t *testing.T) {
	if os.Getenv("ENABLE_INTEGRATION_TESTS") != "1" {
		return
	}
	t.Run("verify save state", func(t *testing.T) {
		sshConnectionParams := newSshConnectionParams()
		sshConnectionParams.initPrivateKey()
		if host, err := os.LookupEnv("INTEGRATION_TEST_SSH_HOST"); err != false {
			sshConnectionParams.hostname = host
		} else {
			sshConnectionParams.hostname = "localhost"
		}
		if port, err := os.LookupEnv("INTEGRATION_TEST_SSH_PORT"); err != false {
			port, _ := strconv.Atoi(port)
			sshConnectionParams.port = int16(port)
		} else {
			sshConnectionParams.port = 2222
		}
		sshConnectionParams.port = 2222
		sshConnection := &SSHConnnection{params: sshConnectionParams}
		err := sshConnection.Connect()
		assert.Equal(t, err, nil)

		executionEnvironment := make(map[string]string)
		executionEnvironment["VARIABLE_1"] = "*Variable \"1\" value*"
		cmd := "echo \"Hello World! Value = ${VARIABLE_1}\""

		var stdOut strings.Builder
		// stdOutBuffer.Write()
		// stdOut := os.Stdout
		err = sshConnection.Execute(executionEnvironment, cmd, &stdOut, &stdOut)
		assert.Equal(t, err, nil)

		output := stdOut.String()
		expectedOutput := "Hello World! Value = *Variable \"1\" value*\n"
		assert.Equal(t, output, expectedOutput)
	})
}
