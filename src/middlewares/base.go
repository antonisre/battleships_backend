package middlewares

import (
	"net/http"

	"github.com/jinzhu/gorm"
)

// BaseMiddleware struct
type BaseMiddleware struct {
	DB *gorm.DB
}

// SetContentTypeHeader to JSON
func SetContentTypeHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(response, request)
	})
}
