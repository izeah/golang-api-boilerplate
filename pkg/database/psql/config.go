package psql

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

type DBConfig struct {
	Host string
	User string
	Pass string
	Port string
	Name string

	SSLMode string
	TZ      string

	SSHClient *ssh.Client
}

// Config ...
type Config struct {
	Write DBConfig
	Read  []DBConfig
}

var dblogger = logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
	SlowThreshold:             200 * time.Millisecond,
	LogLevel:                  logger.Warn,
	IgnoreRecordNotFoundError: false,
	Colorful:                  true,
})

// DSN ...
func (c Config) DSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", c.Write.Host, c.Write.User, c.Write.Pass, c.Write.Name, c.Write.Port, c.Write.SSLMode, c.Write.TZ)
}

// Open ...
func (c Config) Open() (*gorm.DB, error) {
	dialector := postgres.Open(c.DSN())
	if c.Write.SSHClient != nil {
		driverName := fmt.Sprintf("postgres+ssh+%s+write", c.Write.SSHClient.LocalAddr().String())

		found := false
		for _, d := range sql.Drivers() {
			if d == driverName {
				found = true
			}
		}
		if !found {
			// Now we register the ViaSSHDialer with the ssh connection as a parameter
			sql.Register(driverName, &ViaSSHDialer{Client: c.Write.SSHClient})
		}

		dialector = postgres.New(postgres.Config{
			DriverName: driverName,
			DSN:        c.DSN(),
		})
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger:                 dblogger.LogMode(logger.Info),
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		return nil, err
	}

	if len(c.Read) > 0 {
		var replica []gorm.Dialector

		for i, config := range c.Read {
			dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", config.Host, config.User, config.Pass, config.Name, config.Port, config.SSLMode, config.TZ)
			dialector := postgres.Open(dsn)
			if config.SSHClient != nil {
				driverName := fmt.Sprintf("postgres+ssh+%s+read_%d", config.SSHClient.LocalAddr().String(), i)

				found := false
				for _, d := range sql.Drivers() {
					if d == driverName {
						found = true
					}
				}
				if !found {
					// Now we register the ViaSSHDialer with the ssh connection as a parameter
					sql.Register(driverName, &ViaSSHDialer{Client: config.SSHClient})
				}

				dialector = postgres.New(postgres.Config{
					DriverName: driverName,
					DSN:        dsn,
				})
			}
			replica = append(replica, dialector)
		}

		if err = db.Use(dbresolver.Register(dbresolver.Config{
			Replicas:          replica,
			TraceResolverMode: true,
		}).
			SetConnMaxIdleTime(time.Hour).       // SetConnMaxIdleTime sets the maximum amount of time a connection may be idle.
			SetConnMaxLifetime(5 * time.Minute). // SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
			SetMaxIdleConns(5).                  // SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
			SetMaxOpenConns(10),
		); err != nil {
			return nil, err
		}
	}
	return db, nil
}
