package db

import (
	"fmt"
	"github.com/cloudfoundry/uaa-key-rotator/config"
	"strconv"
)

func ConnectionURI(rotatorConfig *config.RotatorConfig) (string, error) {
	var connStr string
	switch rotatorConfig.DatabaseScheme {
	case "mysql":
		{
			timeout := 120

			port, err := strconv.Atoi(rotatorConfig.DatabasePort)
			if err != nil {
				return "", err
			}
			connStr = fmt.Sprintf(
				"%s:%s@tcp(%s:%d)/%s?parseTime=true&timeout=%ds&readTimeout=%ds&writeTimeout=%ds",
				rotatorConfig.DatabaseUsername,
				rotatorConfig.DatabasePassword,
				rotatorConfig.DatabaseHostname,
				port,
				rotatorConfig.DatabaseName,
				timeout,
				timeout,
				timeout,
			)

			if rotatorConfig.DatabaseTlsEnabled {
				if rotatorConfig.DatabaseSkipSSLValidation {
					connStr += "&tls=skip-verify"
				} else {
					connStr += "&tls=true"
				}
			}
		}
	case "postgres":
		connStr = fmt.Sprintf("%s://%s:%s@%s:%s/%s?connect_timeout=120",
			rotatorConfig.DatabaseScheme,
			rotatorConfig.DatabaseUsername,
			rotatorConfig.DatabasePassword,
			rotatorConfig.DatabaseHostname,
			rotatorConfig.DatabasePort,
			rotatorConfig.DatabaseName,
		)

		if rotatorConfig.DatabaseTlsEnabled {
			if rotatorConfig.DatabaseSkipSSLValidation {
				connStr += "&sslmode=require"
			} else {
				connStr += "&sslmode=verify-ca"
			}
		} else {
			connStr += "&sslmode=disable"
		}

	case "sqlserver":
		connStr = fmt.Sprintf("%s://%s:%s@%s:%s?database=%s&connection+timeout=120",
			rotatorConfig.DatabaseScheme,
			rotatorConfig.DatabaseUsername,
			rotatorConfig.DatabasePassword,
			rotatorConfig.DatabaseHostname,
			rotatorConfig.DatabasePort,
			rotatorConfig.DatabaseName,
		)

		if rotatorConfig.DatabaseTlsEnabled {
			if rotatorConfig.DatabaseSkipSSLValidation {
				connStr += "&TrustServerCertificate=true"
			} else {
				connStr += "&TrustServerCertificate=false"
			}
		}
	}
	return connStr, nil
}
