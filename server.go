package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"

	"github.com/codequest-eu/burnafterreading/internal/authorizer"
	"github.com/codequest-eu/burnafterreading/internal/storage"
	"github.com/codequest-eu/burnafterreading/lib"
)

func getS3Storage() (lib.Storage, error) {
	log.Printf("Using S3 storage backend")
	return storage.S3Storage(
		os.Getenv("AWS_ACCESS_KEY_ID"),
		os.Getenv("AWS_SECRET_ACCESS_KEY"),
		os.Getenv("AWS_REGION"),
		os.Getenv("AWS_BUCKET"),
	)
}

func getFileStorage() (lib.Storage, error) {
	log.Printf("Using local file storage backend")
	return storage.LocalFileStorage(os.Getenv("BASE_PATH"))
}

func getStorage() (lib.Storage, error) {
	if os.Getenv("BASE_PATH") != "" {
		return getFileStorage()
	}
	return getS3Storage()
}

func getAuth() lib.Authorizer {
	log.Printf("Using Basic HTTP authorization")
	return authorizer.BasicHTTPAuthorizer(os.Getenv("USER"), os.Getenv("PASS"))
}

func runServer() error {
	addr := os.Getenv("ADDR")
	if os.Getenv("USE_SSL") == "" {
		log.Printf("Starting HTTP server bound to %q", addr)
		return http.ListenAndServe(addr, nil)
	}
	log.Printf("Starting HTTPS server bound to %q", addr)
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
