package database

import (
	"fmt"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// SSHConfig ...
type SSHConfig struct {
	Host string
	Port string
	User string
	Pass string
}

// Open ...
func (c SSHConfig) Open() (Net net.Conn, SSH *ssh.Client, err error) {
	defer func() {
		if err != nil {
			if SSH != nil {
				_ = SSH.Close()
			}
			if Net != nil {
				_ = Net.Close()
			}
		}
	}()

	var agentClient agent.Agent

	// Establish a connection to the local ssh-agent
	if Net, err = net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		// Create a new instance of the ssh agent
		agentClient = agent.NewClient(Net)
	}

	// The client configuration with configuration option to use the ssh-agent
	sshConfig := &ssh.ClientConfig{
		User:            c.User,
		Auth:            []ssh.AuthMethod{},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// When the agentClient connection succeeded, add them as AuthMethod
	if agentClient != nil {
		sshConfig.Auth = append(sshConfig.Auth, ssh.PublicKeysCallback(agentClient.Signers))
	}

	// When there's a non empty password add the password AuthMethod
	sshConfig.Auth = append(sshConfig.Auth, ssh.PasswordCallback(func() (string, error) {
		return c.Pass, nil
	}))

	// Connect to the SSH Server
	SSH, err = ssh.Dial("tcp", fmt.Sprintf("%s:%s", c.Host, c.Port), sshConfig)
	return
}
