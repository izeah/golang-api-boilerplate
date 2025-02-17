package psql

import (
	"net"
	"time"

	"database/sql/driver"
	"github.com/lib/pq"
	"golang.org/x/crypto/ssh"
)

// ViaSSHDialer ...
type ViaSSHDialer struct {
	Client *ssh.Client
}

// Open ...
func (dialer *ViaSSHDialer) Open(s string) (_ driver.Conn, err error) {
	return pq.DialOpen(dialer, s)
}

// Dial ...
func (dialer *ViaSSHDialer) Dial(network, address string) (net.Conn, error) {
	return dialer.Client.Dial(network, address)
}

// DialTimeout ...
func (dialer *ViaSSHDialer) DialTimeout(network, address string, _ time.Duration) (net.Conn, error) {
	return dialer.Client.Dial(network, address)
}
