package main

import (
	"crypto/tls"
	"net/http"
	"os"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/codequest-eu/burnafterreading/internal/authorizer"
	"github.com/codequest-eu/burnafterreading/internal/storage"
	"github.com/codequest-eu/burnafterreading/lib"
	"github.com/polds/logrus-papertrail-hook"
)

func getS3Storage() (lib.Storage, error) {
	log.Infof("Using S3 storage backend")
	return storage.S3Storage(
		os.Getenv("AWS_ACCESS_KEY_ID"),
		os.Getenv("AWS_SECRET_ACCESS_KEY"),
		os.Getenv("AWS_REGION"),
		os.Getenv("AWS_BUCKET"),
	)
}

func getFileStorage() (lib.Storage, error) {
	log.Infof("Using local file storage backend")
	return storage.LocalFileStorage(os.Getenv("BASE_PATH"))
}

func getStorage() (lib.Storage, error) {
	if os.Getenv("BASE_PATH") != "" {
		return getFileStorage()
	}
	return getS3Storage()
}

func getAuth() lib.Authorizer {
	log.Infof("Using Basic HTTP authorization")
	return authorizer.BasicHTTPAuthorizer(os.Getenv("USER"), os.Getenv("PASS"))
}

func runServer() error {
	addr := os.Getenv("ADDR")
	if os.Getenv("USE_SSL") == "" {
		log.Infof("Starting HTTP server bound to %q", addr)
		return http.ListenAndServe(addr, nil)
	}
	log.Infof("Starting HTTPS server bound to %q", addr)
	certificate, err := tls.X509KeyPair(
		[]byte(os.Getenv("CERT_PEM")),
		[]byte(os.Getenv("CERT_KEY")),
	)
	if err != nil {
		log.Fatal(err)
	}
	listener, err := tls.Listen(
		"tcp",
		addr,
		&tls.Config{Certificates: []tls.Certificate{certificate}},
	)
	if err != nil {
		return err
	}
	return http.Serve(listener, nil)
}

func main() {
	storageEngine, err := getStorage()
	if err != nil {
		log.Fatalf("error initializing storage: %v", err)
	}
	http.Handle("/", &lib.Handler{
		Authorizer: getAuth(),
		Storage:    storageEngine,
	})
	log.Fatal(runServer())
}

func init() {
	log.SetOutput(os.Stderr)
	log.SetLevel(log.DebugLevel)
	papertrailHost := os.Getenv("PAPERTRAIL_HOST")
	if papertrailHost == "" {
		return
	}
	papertrailPort := os.Getenv("PAPERTRAIL_PORT")
	port, err := strconv.Atoi(papertrailPort)
	if err != nil {
		log.Fatalf("Error getting Papertrail port: %v", err)
	}
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("Error getting hostname: %v", err)
	}
	hook, err := logrus_papertrail.NewPapertrailHook(
		&logrus_papertrail.Hook{
			Host:     papertrailHost,
			Port:     port,
			Hostname: hostname,
			Appname:  os.Getenv("APP_NAME"),
		},
	)
	if err != nil {
		log.Fatalf("Error building Papertrail hook: %v", err)
	}
	log.AddHook(hook)
}
