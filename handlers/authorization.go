package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/energietransitie/twomes-backoffice-api/pkg/ports"
	"github.com/energietransitie/twomes-backoffice-api/pkg/twomes"
)

// A Contextkey is the type for a context key.
type contextKey int

// AuthorizationCtxKey is the key for the authorization value that is passed to the context,
// when the authentication middleware is used.
const AuthorizationCtxKey contextKey = 0

type AuthorizationHandler struct {
	service ports.AuthorizationService
}

// Create a new AuthorizationHandler.
func NewAuthorizationHandler(service ports.AuthorizationService) *AuthorizationHandler {
	return &AuthorizationHandler{
		service: service,
	}
}

func (h *AuthorizationHandler) Middleware(kind twomes.AuthKind) func(next Handler) Handler {
	return func(next Handler) Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				return NewHandlerError(nil, "unauthorized", http.StatusUnauthorized).WithMessage("authorization header not present")
			}

			authHeader = strings.Split(authHeader, "Bearer ")[1]

			if authHeader == "" {
				return NewHandlerError(nil, "unauthorized", http.StatusUnauthorized).WithMessage("authorization malformed")
			}

			auth, err := h.service.ParseTokenToAuthorization(authHeader)
			if err != nil {
				return NewHandlerError(err, "unauthorized", http.StatusUnauthorized).WithMessage("error when parsing token")
			}

			if !auth.IsKind(kind) {
				return NewHandlerError(nil, "unauthorized", http.StatusUnauthorized).WithMessage("incorrect authorization kind was used to access route")
			}

			// Add the value of audience to the HTTP context with key AuthenticatedID.
			authCtx := context.WithValue(r.Context(), AuthorizationCtxKey, auth)
			r = r.WithContext(authCtx)

			return next(w, r)
		}
	}
}
