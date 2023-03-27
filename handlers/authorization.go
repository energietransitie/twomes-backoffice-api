package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/energietransitie/twomes-api/pkg/ports"
	"github.com/energietransitie/twomes-api/pkg/twomes"
	"github.com/sirupsen/logrus"
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

func (h *AuthorizationHandler) Middleware(kind twomes.AuthKind) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				logrus.Info("authorization header not present")
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			authHeader = strings.Split(authHeader, "Bearer ")[1]

			if authHeader == "" {
				logrus.Info("authorization malformed")
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			auth, err := h.service.ParseTokenToAuthorization(authHeader)
			if err != nil {
				logrus.WithField("error", err).Info("error when parsing token")
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			if !auth.IsKind(kind) {
				logrus.WithFields(logrus.Fields{
					"route":        r.URL.Path,
					"kindProvided": auth.Kind,
					"kindNeeded":   kind,
				}).Info("incorrect authorization kind was used to access route")
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			// Add the value of audience to the HTTP context with key AuthenticatedID.
			authCtx := context.WithValue(r.Context(), AuthorizationCtxKey, auth)
			r = r.WithContext(authCtx)

			next.ServeHTTP(w, r)
		}
	}
}
