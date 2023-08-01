package api

import (
	"net/http"

	"github.com/praveenmahasena647/server/internal/helpers"
)

func StartServer() error {
	http.HandleFunc("/", helpers.Serve)
	return http.ListenAndServe(":42069", nil)
}
