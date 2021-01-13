package middlewares

import (
	"fmt"
	"net/http"
	"time"

	"github.com/charliemcelfresh/kata/internal/config"
)

func EnforceAPIKataRequestContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")

		if contentType != config.Constants["REQUIRED_API_KATA_REQUEST_CONTENT_TYPE"] {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			fmt.Fprintln(w, fmt.Sprintf("Content-Type header must be %v",
				config.Constants["REQUIRED_API_KATA_REQUEST_CONTENT_TYPE"]),
			)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func AddResponseContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()
		header.Set("Content-Type", config.Constants["API_KATA_RESPONSE_CONTENT_TYPE"].(string))
		next.ServeHTTP(w, r)
	})
}

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(fmt.Sprintf(
			`timestamp: %v, request_path: %v`,
			time.RFC3339, r.RequestURI,
		))
		next.ServeHTTP(w, r)
	})
}
