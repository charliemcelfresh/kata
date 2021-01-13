package cmd

import (
	"fmt"
	"net/http"

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
	mux := http.NewServeMux()
	rootHandler := http.HandlerFunc(rootHandler)
	sharedMiddlewares := alice.New(
		middlewares.AddResponseContentType,
		middlewares.EnforceAPIKataRequestContentType,
		middlewares.LogRequest,
	)
	mux.Handle("/", sharedMiddlewares.Then(rootHandler))
	logrus.Info(fmt.Sprintf("kata server running on %v",
		config.Constants["SERVER_PORT"]))

	err := http.ListenAndServe(fmt.Sprintf("%v", config.Constants["SERVER_PORT"]), mux)
	if err != nil {
		logrus.Fatal(err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{"Hello": "Kata!"}`))
}
