// http.go contains useful http utilities.
package utils

import (
	"net/http"

	"github.com/docker/vetinari/auth"
	"github.com/docker/vetinari/errors"
)

type BetterHandler func(ctx IContext, w http.ResponseWriter, r *http.Request) *errors.DockerError

func errorHandler(handler BetterHandler) {
	errorWrapper := func(w http.ResponseWriter, r *http.Request) {
		if err := handler(); err != nil {
			// TODO: Log error
			http.Error(w, err.Error(), err.HTTPStatus)
		}
	}
	return errorWrapper
}

func BaseHandler(handler BetterHandler) http.Handler {
	baseWrapper := func(w http.ResponseWriter, r *http.Request) *errors.DockerError {
		ctx := generateContext(r)
		return handler(ctx, w, r)
	}
	return errorHandler(baseWrapper)
}

func AuthorizedHandler(handler BetterHandler, auth IAuthorizer, scopes ...Scope) http.Handler {
	authorizedWrapper := func(ctx IContext, w http.ResponseWriter, r *http.Request) errors.DockerError {
		if err := auth.Authorize(ctx, scopes...); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		return handler(ctx, w, r)
	}
	return BaseHandler(authorizedWrapper)
}