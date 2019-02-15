package cli

import (
	"net/http"

	"github.com/lucab/exp-locksmith2/internal/daemon"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	cmdDaemon = &cobra.Command{
		Use:  "daemon",
		RunE: runDaemon,
	}
	addr = "127.0.0.1:9999"
)

func init() {
	locksmith2Cmd.AddCommand(cmdDaemon)
}

func runDaemon(cmd *cobra.Command, cmdArgs []string) error {
	logrus.Info("starting daemon on ", addr)

	http.HandleFunc("/v1/pre-reboot", daemon.PreReboot)
	http.HandleFunc("/v1/steady-state", daemon.SteadyState)
	return http.ListenAndServe(addr, nil)
}
