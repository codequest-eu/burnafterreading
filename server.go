package main

import (
	"log"
	"net/http"
	"os"

	"github.com/codequest-eu/burnafterreading/internal/authorizer"
	"github.com/codequest-eu/burnafterreading/internal/storage"
	"github.com/codequest-eu/burnafterreading/lib"
)

func getS3Storage() (lib.Storage, error) {
	return storage.S3Storage(
		os.Getenv("AWS_ACCESS_KEY_ID"),
		os.Getenv("AWS_SECRET_ACCESS_KEY"),
		os.Getenv("AWS_REGION"),
		os.Getenv("AWS_BUCKET"),
	)
}

func getFileStorage() (lib.Storage, error) {
	return storage.LocalFileStorage(os.Getenv("BASE_PATH"))
}

func getStorage() (lib.Storage, error) {
	if os.Getenv("BASE_PATH") != "" {
		return getFileStorage()
	}
	return getS3Storage()
}

func getAuth() lib.Authorizer {
	return authorizer.BasicHTTPAuthorizer(os.Getenv("USER"), os.Getenv("PASS"))
}

func runServer() error {
	addr := os.Getenv("ADDR")
	log.Printf("Starting HTTP server bound to %q", addr)
	if os.Getenv("USE_SSL") == "" {
		return http.ListenAndServe(addr, nil)
	}
	return http.ListenAndServeTLS(addr, "cert.pem", "cert.key", nil)
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
