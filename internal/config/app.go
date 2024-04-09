package config

import (
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"

	"boilerplate/pkg/util/priority"
)

type AppConfig struct {
	ENV     string
	Name    string
	Host    string
	Port    int
	Version string
	Schemes []string
}

var (
	appConfig *AppConfig
	appOnce   sync.Once
)

func App() *AppConfig {
	appOnce.Do(func() {
		appConfig = new(AppConfig)

		appConfig.Version = "1.0.0"
		appConfig.Name = priority.PriorityString(os.Getenv("APP_NAME"), "baf-qrcode-api-cms")
		appConfig.ENV = strings.ToLower(strings.TrimSpace(os.Getenv("ENV")))
		if appConfig.ENV != "development" && appConfig.ENV != "staging" && appConfig.ENV != "production" {
			appConfig.ENV = "local"
		}

		port := 4040
		var err error
		if strPort, isExist := os.LookupEnv("PORT"); isExist {
			if port, err = strconv.Atoi(strPort); err != nil {
				panic(err)
			}
		}
		appConfig.Port = port
		appConfig.Host = os.Getenv("HOST")

		scheme := strings.ToLower(strings.TrimSpace(os.Getenv("SCHEMES")))
		splitedSchemes := strings.Split(scheme, ",")
		appSchemes := slices.Compact(splitedSchemes)
		for _, appScheme := range appSchemes {
			if appScheme == "http" || appScheme == "https" {
				appConfig.Schemes = append(appConfig.Schemes, appScheme)
			}
		}

		if len(appConfig.Schemes) < 1 {
			appConfig.Schemes = []string{"http"}
		}
	})
	return appConfig
}

func (c *AppConfig) IsLocal() bool {
	return c.ENV == "local"
}
