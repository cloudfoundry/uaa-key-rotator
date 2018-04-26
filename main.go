package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/cloudfoundry/uaa-key-rotator/config"
	"github.com/cloudfoundry/uaa-key-rotator/crypto"
	db2 "github.com/cloudfoundry/uaa-key-rotator/db"
	"github.com/cloudfoundry/uaa-key-rotator/rotator"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGTERM, os.Interrupt)

	fmt.Println("rotator has started")

	configPath := flag.String("config", "", "Path to uaa key rotator config file")
	flag.Parse()
	configFile, err := os.Open(*configPath)
	if err != nil {
		panic(err)
	}

	rotatorConfig, err := config.New(configFile)
	if err != nil {
		panic(err)
	}

	var rotatorChan = make(chan struct{})
	go rotate(rotatorConfig, rotatorChan)

	select {
	case <-rotatorChan:
		os.Exit(0)
	case s := <-sigChan:
		if s == os.Interrupt {
			fmt.Println("shutting down gracefully...")
			os.Exit(1)
		}
	}
}

func rotate(rotatorConfig *config.RotatorConfig, rotatorChan chan struct{}) {
	connStr := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable",
		rotatorConfig.DatabaseScheme,
		rotatorConfig.DatabaseUsername,
		rotatorConfig.DatabasePassword,
		rotatorConfig.DatabaseHostname,
		rotatorConfig.DatabasePort,
		rotatorConfig.DatabaseName,
	)
	db, err := getDbConn(connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	credentials, err := db2.ReadAll(db)
	if err != nil {
		panic(err)
	}
	keyService := rotator.UaaKeyService{
		ActiveKeyLabel: rotatorConfig.ActiveKeyLabel,
		EncryptionKeys: rotatorConfig.EncryptionKeys,
	}
	r := rotator.UAARotator{
		KeyService:     keyService,
		SaltAccessor:   crypto.UaaSaltAccessor{},
		NonceAccessor:  crypto.UaaNonceAccessor{},
		CipherAccessor: crypto.UAACipherAccessor{},
		DbMapper:       rotator.DbMapper{},
	}
	for _, cred := range credentials {
		rotatedCred, err := r.Rotate(cred)
		if err != nil {
			panic(err)
		}
		if err = db2.Write(db, rotatedCred); err != nil {
			panic(err)
		}
	}
	fmt.Println("rotator has finished")
	close(rotatorChan)
}

func getDbConn(connectionString string) (*sqlx.DB, error) {
	nativeDBConn, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("unable to open database connection: %s", err)
	}

	dbConn := sqlx.NewDb(nativeDBConn, "postgres")

	if err = dbConn.Ping(); err != nil {
		dbConn.Close()
		if netErr, ok := err.(*net.OpError); ok {
			return nil, errors.Wrap(netErr, "unable to ping")
		}
		return nil, fmt.Errorf("unable to ping: %s", err)
	}

	return dbConn, nil
}

//TODO: tls certs
/*
db, err := sql.Open("mysql", "user@tcp(localhost:3306)/test?tls=custom")
    db.Driver().(mysql.TLSConfig).SetTLSConfig("custom", &tls.Config{
        RootCAs:      rootCAs,
        Certificates: clientCerts,
    })
*/
