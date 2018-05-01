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
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"net"
	"os"
	"os/signal"
	"syscall"
	_ "code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/lager"
	"sync"
	"context"
	"runtime"
	"log"
)

func main() {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGTERM, os.Interrupt)

	allowThreadDumpOnSigQUIT()

	logger := lager.NewLogger(fmt.Sprintf("%s.%s", "rotator", "uaa-key-rotator"))
	logger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.INFO))

	logger.Info("rotator has started")

	configPath := flag.String("config", "", "Path to uaa key rotator config file")
	flag.Parse()
	configFile, err := os.Open(*configPath)
	if err != nil {
		logger.Fatal("unable to open config", err)
	}

	rotatorConfig, err := config.New(configFile)
	if err != nil {
		logger.Fatal("unable to parse config", err)
	}

	var rotatorChan = make(chan struct{})
	parentCtx := context.Background()
	rotatorCtx, cancelRotatorFunc := context.WithCancel(parentCtx)
	go rotate(rotatorCtx, logger, rotatorConfig, rotatorChan)

	select {
	case s := <-sigChan:
		if s == os.Interrupt {
			logger.Info("shutting down gracefully...")
			cancelRotatorFunc()
		}
	case <-rotatorChan:
		os.Exit(0)
	}
}

func rotate(parentCtx context.Context, logger lager.Logger, rotatorConfig *config.RotatorConfig, rotatorChan chan struct{}) {
	defer close(rotatorChan)

	db, err := getDbConn(rotatorConfig.DatabaseScheme, db2.ConnectionURI(rotatorConfig))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	credentialsDBFetcher := db2.GoogleMfaCredentialsDBFetcher{
		DB:             db,
		ActiveKeyLabel: rotatorConfig.ActiveKeyLabel,
	}

	var fetcherErrChan <-chan error
	credentialsChan, fetcherErrChan := credentialsDBFetcher.RowsToRotate()
	credentialsDBUpdater := db2.GoogleMfaCredentialsDBUpdater{
		DB: db,
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

	ctx, cancel := context.WithCancel(parentCtx)

	worker := func(wg *sync.WaitGroup) {
		defer wg.Done()

		for {
			select {
			case cred, ok := <-credentialsChan:
				if !ok {
					logger.Debug("No more mfa credentials. Worker signing off...")
					return
				}

				logger.Info("rotating mfa cred", lager.Data{"mfa_cred": cred})
				rotatedCred, err := r.Rotate(cred)
				if err != nil {
					logger.Error("unable to rotate record... Skipping", err)
					continue
				}

				err = credentialsDBUpdater.Write(rotatedCred)
				if err != nil {
					logger.Error("unable to update record... Skipping", err)
					continue
				}

			case err := <-fetcherErrChan:
				logger.Error("error during fetching a record...", err)
				cancel()
			case <-ctx.Done():
				logger.Info("rotator worker has been cancelled")
				return
			}
		}
	}

	wg := sync.WaitGroup{}

	numWorkers := 4
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go worker(&wg)
	}

	logger.Info("workers are unleahsed")
	wg.Wait()
	logger.Info("rotator has finished")
}

func getDbConn(scheme string, connectionString string) (db2.Queryer, error) {
	nativeDBConn, err := sql.Open(scheme, connectionString)
	if err != nil {
		return nil, fmt.Errorf("unable to open database connection: %s", err)
	}

	dbConn := sqlx.NewDb(nativeDBConn, scheme)
	if err = dbConn.Ping(); err != nil {
		dbConn.Close()
		if netErr, ok := err.(*net.OpError); ok {
			return nil, errors.Wrap(netErr, "unable to ping")
		}
		return nil, fmt.Errorf("unable to ping: %s", err)
	}

	return db2.DbAwareQuerier{DB: dbConn, DBScheme: scheme}, nil
}

func allowThreadDumpOnSigQUIT() {
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGQUIT)
		buf := make([]byte, 1<<20)
		for {
			<-sigs
			stacklen := runtime.Stack(buf, true)
			log.Printf("=== received SIGQUIT ===\n*** goroutine dump...\n%s\n*** end\n", buf[:stacklen])
		}
	}()
}