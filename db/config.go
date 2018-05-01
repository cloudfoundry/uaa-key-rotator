package db

import (
	"github.com/cloudfoundry/uaa-key-rotator/config"
	"strconv"
	"fmt"
)

func ConnectionURI(rotatorConfig *config.RotatorConfig) string {
	var connStr string
	switch rotatorConfig.DatabaseScheme {
	case "mysql":
		{
			timeout := 120

			port, err := strconv.Atoi(rotatorConfig.DatabasePort)
			if err != nil {
				panic(err)
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
				connStr += "&useSSL=true"
				connStr += fmt.Sprintf("&trustServerCertificate=%t", rotatorConfig.DatabaseSkipSSLValidation)
				if len(rotatorConfig.DatabaseTLSProtocols) != 0 {
					connStr += "&enabledSslProtocolSuites=" + rotatorConfig.DatabaseTLSProtocols
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
	}
	return connStr
}
