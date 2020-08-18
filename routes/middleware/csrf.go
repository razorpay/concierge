package middleware

import (
	"concierge/config"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
	adapter "github.com/gwatts/gin-adapter"
)

var csrfMd func(http.Handler) http.Handler

func init() {
	csrfMd = csrf.Protect([]byte(config.CSRFConfig.AuthKey),
		csrf.MaxAge(0),
		csrf.Secure(config.CSRFConfig.Secure),
		csrf.FieldName("_token"),
		csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"message": "Forbidden - CSRF token invalid"}`))
		})),
	)
}

// CSRF ...
func CSRF() gin.HandlerFunc {
	return adapter.Wrap(csrfMd)
}
