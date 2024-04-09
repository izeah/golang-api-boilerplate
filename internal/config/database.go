package config

import (
	"os"
	"strings"
	"sync"

	"boilerplate/pkg/util/priority"
)

type DBConfig struct {
	Host    string
	Port    string
	Name    string
	User    string
	Pass    string
	SSLMode string
	TZ      string
}

var (
	db     *DBConfig
	dbOnce sync.Once
)

func DB() *DBConfig {
	dbOnce.Do(func() {
		listSSLMode := map[string]bool{
			"disable":     true,
			"allow":       true,
			"prefer":      true,
			"require":     true,
			"verify-ca":   true,
			"verify-full": true,
		}
		givenSSLMode := strings.ToLower(strings.TrimSpace(priority.PriorityString(os.Getenv("DB_SSLMODE"), "disable")))
		if _, ok := listSSLMode[givenSSLMode]; !ok {
			givenSSLMode = "disable"
		}

		db = &DBConfig{
			Host:    priority.PriorityString(os.Getenv("DB_HOST"), "localhost"),
			Port:    priority.PriorityString(os.Getenv("DB_PORT"), "5432"),
			Name:    priority.PriorityString(os.Getenv("DB_NAME"), "db_name"),
			User:    priority.PriorityString(os.Getenv("DB_USER"), "postgres"),
			Pass:    priority.PriorityString(os.Getenv("DB_PASS"), "postgres"),
			TZ:      priority.PriorityString(os.Getenv("DB_TZ"), "UTC"),
			SSLMode: givenSSLMode,
		}
	})
	return db
}
