package database

import (
	"database/sql"
	"fmt"
	"net"
	"os"
	"time"

	"boilerplate/internal/config"
	"boilerplate/pkg/database/msql"
	"boilerplate/pkg/database/psql"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"gorm.io/gorm"
)

var (
	err           error
	netConnection map[string]net.Conn
	sshConnection map[string]*ssh.Client
	dbConnections map[string]*gorm.DB
)

// Database ...
type Database interface {
	Open() (*gorm.DB, error)
	DSN() string
}

// Init ...
func Init() {
	netConnection = make(map[string]net.Conn)
	sshConnection = make(map[string]*ssh.Client)
	sshConfig := map[string]SSHConfig{
		//"PGSQL": {
		//	Host: os.Getenv("PGSQL_SSH_HOST"),
		//	User: os.Getenv("PGSQL_SSH_USER"),
		//	Pass: os.Getenv("PGSQL_SSH_PASS"),
		//	Port: os.Getenv("PGSQL_SSH_PORT"),
		//},
		//"PGSQL_REPLICA_1": {
		//	Host: os.Getenv("PGSQL_SSH_REPLICA_1_HOST"),
		//	User: os.Getenv("PGSQL_SSH_REPLICA_1_USER"),
		//	Pass: os.Getenv("PGSQL_SSH_REPLICA_1_PASS"),
		//	Port: os.Getenv("PGSQL_SSH_REPLICA_1_PORT"),
		//},
		"MYSQL": {
			Host: os.Getenv("MYSQL_SSH_HOST"),
			Port: os.Getenv("MYSQL_SSH_PORT"),
			User: os.Getenv("MYSQL_SSH_USER"),
			Pass: os.Getenv("MYSQL_SSH_PASS"),
		},
	}
	for k, config := range sshConfig {
		if netConnection[k], sshConnection[k], err = config.Open(); err != nil {
			panic(fmt.Errorf("connection to ssh %s, error: %v", k, err))
		}
		logrus.Info(fmt.Sprintf("successfully connected to ssh %s", k))
	}

	dbConnections = make(map[string]*gorm.DB)
	dbConfig := map[string]Database{
		"PGSQL_DB_BAF": psql.Config{
			Write: psql.DBConfig{
				Host:    config.DB().Host,
				User:    config.DB().User,
				Pass:    config.DB().Pass,
				Port:    config.DB().Port,
				Name:    config.DB().Name,
				SSLMode: config.DB().SSLMode,
				TZ:      config.DB().TZ,
				// SSHClient: sshConnection["PGSQL"],
			},
			// Read: []psql.DBConfig{
			// {
			// Host:    os.Getenv("PGSQL_DB_BAF_REPLICA_1_HOST"),
			// User:    os.Getenv("PGSQL_DB_BAF_REPLICA_1_USER"),
			// Pass:    os.Getenv("PGSQL_DB_BAF_REPLICA_1_PASS"),
			// Port:    os.Getenv("PGSQL_DB_BAF_REPLICA_1_PORT"),
			// Name:    os.Getenv("PGSQL_DB_BAF_REPLICA_1_NAME"),
			// SslMode: os.Getenv("PGSQL_DB_BAF_REPLICA_1_SSLMODE"),
			// Tz:      os.Getenv("PGSQL_DB_BAF_REPLICA_1_TZ"),
			// SSHClient: sshConnection["PGSQL_REPLICA_1"],
			// },
			// },
		},
		"MYSQL": msql.Config{
			Write: msql.DBConfig{
				Host:      os.Getenv("MYSQL_DB_HOST"),
				User:      os.Getenv("MYSQL_DB_USER"),
				Pass:      os.Getenv("MYSQL_DB_PASS"),
				Port:      os.Getenv("MYSQL_DB_PORT"),
				Name:      os.Getenv("MYSQL_DB_NAME"),
				ParseTime: os.Getenv("MYSQL_DB_PARSE_TIME"),
				SSHClient: sshConnection["MYSQL"],
			},
		},
	}
	var sqlDB *sql.DB
	for k, db := range dbConfig {
		if dbConnections[k], err = db.Open(); err != nil {
			panic(fmt.Errorf("connection to db %s, error: %v", k, err))
		}

		if sqlDB, err = dbConnections[k].DB(); err != nil {
			panic(fmt.Errorf("connection to db %s, error: %v", k, err))
		}

		sqlDB.SetMaxOpenConns(20)
		sqlDB.SetMaxIdleConns(5)
		sqlDB.SetConnMaxLifetime(5 * time.Minute)
		sqlDB.SetConnMaxIdleTime(time.Hour)

		if err = sqlDB.Ping(); err != nil {
			panic(fmt.Errorf("connection to db %s, error: %v", k, err))
		}

		logrus.Info(fmt.Sprintf("successfully connected to db %s", k))
	}
}

// Connection ...
func Connection(name string) *gorm.DB {
	if dbConnections[name] == nil {
		panic("connection is undefined")
	}
	return dbConnections[name]
}

// PSQL ...
func PSQL() *gorm.DB {
	return Connection("PGSQL_DB_BAF")
}

// MYSQL ...
func MYSQL() *gorm.DB {
	return Connection("MYSQL")
}

// Close ...
func Close() {
	var sqlDB *sql.DB
	for k, db := range dbConnections {
		if sqlDB, err = db.DB(); err == nil {
			err = sqlDB.Close()
		}
		if err != nil {
			logrus.WithField("message", "failed to close db connection "+k).Error(err.Error())
		} else {
			logrus.Infof("db connection to %s closed", k)
		}
	}

	for key, conn := range sshConnection {
		if conn != nil {
			if err = conn.Close(); err != nil {
				logrus.WithField("message", "failed to close ssh connection "+key).Error(err.Error())
			} else {
				logrus.Infof("ssh connection to %v closed", key)
			}
		}
	}

	for key, conn := range netConnection {
		if conn != nil {
			if err = conn.Close(); err != nil {
				logrus.WithField("message", "failed to close net connection "+key).Error(err.Error())
			} else {
				logrus.Infof("net connection to %v closed", key)
			}
		}
	}
}
