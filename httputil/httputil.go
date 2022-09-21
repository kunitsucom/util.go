package httputil

import (
	"errors"
	"net/http"
)

func ListenAndServe(server *http.Server) error {
	return listenAndServe(server)
}

type httpServer interface {
	ListenAndServe() error
}

func listenAndServe(server httpServer) error {
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		// nolint: wrapcheck
		return err
	}

	return nil
}
