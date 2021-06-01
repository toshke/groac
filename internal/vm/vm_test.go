package vm

import (
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

// func TestSSHConnect(t *testing.T) {
// 	t.Run("verify save state", func(t *testing.T) {
// 		sshConnectionParams := newSshConnectionParams()
// 		sshConnectionParams.initPrivateKey()
// 		sshConnectionParams.hostname = "localhost"
// 		sshConnectionParams.port = 2222
// 		sshConnection := &SSHConnnection{params: sshConnectionParams}
// 		connection, err := sshConnection.Connect()
// 		assert.Equal(t, err, nil)
// 		assert.NotEqual(t, connection, nil)
// 	})
// }
