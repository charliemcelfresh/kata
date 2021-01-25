package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/charliemcelfresh/kata/internal/middlewares"
	"github.com/google/uuid"
	"github.com/gorilla/securecookie"
	"github.com/jmoiron/sqlx"

	"github.com/charliemcelfresh/kata/internal/config"
	_ "github.com/charliemcelfresh/kata/internal/config"
	"github.com/justinas/alice"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run the kata server",
	RunE: func(cmd *cobra.Command, args []string) error {
		serverCmdRunner()
		return nil
	},
}

type server struct {
	logger          middlewares.Logger
	port            string
	shutdownChannel chan os.Signal
	db              *sqlx.DB
	cookieSecurer   *securecookie.SecureCookie
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

func serverCmdRunner() {
	shutdownChannel := make(chan os.Signal, 1)
	signal.Notify(shutdownChannel, os.Interrupt)
	db, err := setupDBConnection()
	if err != nil {
		logrus.WithError(err).Fatal("DB connection failed")
	}
	server := newServer(logrus.New(), config.Constants["SERVER_PORT"].(string), shutdownChannel, db)
	server.serve()
}

func setupDBConnection() (*sqlx.DB, error) {
	return sqlx.Open("postgres", config.Constants["DATABASE_URL"].(string))
}

func newServer(logger middlewares.Logger, port string, shutdownChannel chan os.Signal, db *sqlx.DB) server {
	cookieSecurer := securecookie.New(
		[]byte(config.Constants["COOKIE_HASH_KEY"].(string)),
		[]byte(config.Constants["COOKIE_BLOCK_KEY"].(string)),
	)
	return server{
		logger:          logger,
		port:            port,
		shutdownChannel: shutdownChannel,
		db:              db,
		cookieSecurer:   cookieSecurer,
	}
}

func (s server) serve() {
	mux := http.NewServeMux()
	rootHandler := http.HandlerFunc(s.rootHandler)
	loginHandler := http.HandlerFunc(s.loginHandler)
	middleware := middlewares.NewMiddlewareRunner(s.logger)
	sharedMiddlewares := alice.New(
		middleware.AddResponseContentType,
		middleware.EnforceAPIKataRequestContentType,
		middleware.LogRequest,
	)
	mux.Handle("/", sharedMiddlewares.Then(rootHandler))
	mux.Handle("/login", sharedMiddlewares.Then(loginHandler))
	logrus.Info(fmt.Sprintf("kata server running on %v",
		s.port))

	server := http.Server{Addr: s.port, Handler: mux}
	go server.ListenAndServe()

	<-s.shutdownChannel
	server.Shutdown(context.Background())
}

func (s server) rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{"Hello": "Kata!"}`))
}

func (s server) loginHandler(w http.ResponseWriter, r *http.Request) {
	cookieValue := uuid.New().String()
	_, err := s.db.Exec(`INSERT INTO session (value) VALUES ($1);`, cookieValue)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if encodedValue, err := s.cookieSecurer.Encode("kata-cookie", cookieValue); err == nil {
		cookie := &http.Cookie{
			Name:     "kata-cookie",
			Value:    encodedValue,
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte(`{"Login": "Success"}`))
}
