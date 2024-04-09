package msql

import (
	"context"
	"database/sql/driver"
	"net"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/ssh"
)

// ViaSSHDialer ...
type ViaSSHDialer struct {
	Client *ssh.Client
}

// Open ...
func (dialer *ViaSSHDialer) Open(dsn string) (_ driver.Conn, err error) {
	var cfg *mysql.Config
	if cfg, err = mysql.ParseDSN(dsn); err != nil {
		return nil, err
	}

	mysql.RegisterDialContext(cfg.Net, dialer.Dial)

	var c driver.Connector
	if c, err = mysql.NewConnector(cfg); err != nil {
		return nil, err
	}
	return c.Connect(context.Background())
}

// Dial ...
func (dialer *ViaSSHDialer) Dial(_ context.Context, addr string) (net.Conn, error) {
	return dialer.Client.Dial("tcp", addr)
}
