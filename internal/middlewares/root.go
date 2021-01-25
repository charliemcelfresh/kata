package middlewares

import (
	"fmt"
	"net/http"
	"time"

	"github.com/charliemcelfresh/kata/internal/repository"

	"github.com/google/uuid"

	"github.com/charliemcelfresh/kata/internal/config"
)

var (
	cookieDuration = 60 * 24 * time.Minute
)

type Logger interface {
	Info(...interface{})
}

type Repo interface {
	InsertSessionCookie(value string) (err error)
}

type MiddlewareRunner struct {
	logger     Logger
	repository Repo
}

func NewMiddlewareRunner(logger Logger) MiddlewareRunner {
	r := repository.NewRepository()
	return MiddlewareRunner{
		logger:     logger,
		repository: r,
	}
}

func (m MiddlewareRunner) EnforceAPIKataRequestContentType(next http.Handler) http.Handler {
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

func (m MiddlewareRunner) AddResponseContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()
		header.Set("Content-Type", config.Constants["API_KATA_RESPONSE_CONTENT_TYPE"].(string))
		next.ServeHTTP(w, r)
	})
}

func (m MiddlewareRunner) LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.logger.Info(fmt.Sprintf(
			`timestamp: %v, request_path: %v`,
			time.Now(), r.RequestURI,
		))
		next.ServeHTTP(w, r)
	})
}

func (m MiddlewareRunner) SetCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookieValue := randomString()
		err := m.repository.InsertSessionCookie(cookieValue)
		if err != nil {
			m.logger.Info(err)
		}
		cookie := http.Cookie{
			Name:     "kata-session",
			Value:    cookieValue,
			Path:     "/",
			Domain:   "",
			Expires:  time.Now().Add(cookieDuration),
			MaxAge:   0,
			Secure:   true,
			HttpOnly: true,
			SameSite: 0,
			Raw:      "",
			Unparsed: nil,
		}
		http.SetCookie(w, &cookie)
		next.ServeHTTP(w, r)
	})
}

func randomString() string {
	return uuid.New().String()
}
