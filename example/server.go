package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/codequest-eu/burnafterreading"
	"github.com/codequest-eu/burnafterreading/authorizer"
	"github.com/codequest-eu/burnafterreading/storage"
)

var (
	path = flag.String("path", "/tmp", "path to the directory")
	bind = flag.String("bind", ":1983", "(ip):port to bind to")
	user = flag.String("user", "bacon", "user to use for HTTP auth")
	pass = flag.String("pass", "cabbage", "pass to use for HTTP auth")
)

func main() {
	flag.Parse()
	storage, err := storage.LocalFileStorage(*path)
	if err != nil {
		log.Fatalf("error initializing storage: %v", err)
	}
	http.Handle("/", &burnafterreading.Handler{
		Authorizer: authorizer.BasicHTTPAuthorizer(*user, *pass),
		Storage:    storage,
	})
	log.Fatal(http.ListenAndServe(*bind, nil))
}
