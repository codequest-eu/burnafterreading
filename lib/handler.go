package lib

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

var (
	keyName     = "key"
	errNotFound = errors.New("Not found")
)

// Handler implements a http.Handler interface so that it can be mounted in an
// HTTP server.
type Handler struct {
	Authorizer   Authorizer
	Storage      Storage
	ErrorHandler func(error)
}

func (h *Handler) handlePUT(w http.ResponseWriter, r *http.Request, key string) error {
	writer, err := h.Storage.Put(key)
	if err != nil {
		return err
	}
	defer writer.Close()
	defer r.Body.Close()
	_, err = io.Copy(writer, r.Body)
	return err
}

func (h *Handler) handleGET(w http.ResponseWriter, key string) error {
	reader, err := h.Storage.Get(key)
	if err != nil {
		return err
	}
	defer reader.Close()
	_, err = io.Copy(w, reader)
	return err
}

func (h *Handler) handleDELETE(w http.ResponseWriter, key string) error {
	return h.Storage.Delete(key)
}

func (h *Handler) serve(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	key := r.Form.Get("key")
	if key == "" {
		return errNotFound
	}
	log.Printf("[%s] with key %q from %q", r.Method, key, r.RemoteAddr)
	key = keyAsHash(key)
	if r.Method == "PUT" {
		return h.handlePUT(w, r, key)
	}
	if r.Method == "GET" {
		return h.handleGET(w, key)
	}
	if r.Method == "DELETE" {
		return h.handleDELETE(w, key)
	}
	return errNotFound
}

func keyAsHash(key string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(key)))
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.Authorizer != nil && !h.Authorizer.Authorize(r) {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}
	err := h.serve(w, r)
	if err == nil {
		return
	}
	log.Printf("Error: %v", err)
	w.WriteHeader(http.StatusInternalServerError)
	if h.ErrorHandler != nil {
		h.ErrorHandler(err)
	}
}
