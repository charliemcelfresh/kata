package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/charliemcelfresh/kata/internal/middlewares"

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

func init() {
	rootCmd.AddCommand(serverCmd)
}

func serverCmdRunner() {
	shutdownChannel := make(chan os.Signal, 1)
	signal.Notify(shutdownChannel, os.Interrupt)
	serve(logrus.New(), config.Constants["SERVER_PORT"].(string), shutdownChannel)
}

func serve(logger middlewares.Logger, port string, shutdownChannel chan os.Signal) {
	mux := http.NewServeMux()
	rootHandler := http.HandlerFunc(rootHandler)
	middleware := middlewares.NewMiddlewareRunner(logger)
	sharedMiddlewares := alice.New(
		middleware.AddResponseContentType,
		middleware.EnforceAPIKataRequestContentType,
		middleware.LogRequest,
	)
	mux.Handle("/", sharedMiddlewares.Then(rootHandler))
	logrus.Info(fmt.Sprintf("kata server running on %v",
		port))

	server := http.Server{Addr: port, Handler: mux}
	go server.ListenAndServe()

	<-shutdownChannel
	server.Shutdown(context.Background())
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{"Hello": "Kata!"}`))
}
