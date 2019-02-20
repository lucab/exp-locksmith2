package cli

import (
	"fmt"
	"net/http"
	"time"

	"github.com/lucab/exp-locksmith2/internal/daemon"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	cmdDaemon = &cobra.Command{
		Use:  "daemon",
		RunE: runDaemon,
	}
	address        = "0.0.0.0"
	port           = 9999
	etcdURLs       = []string{"http://127.0.0.1:2379"}
	lockTimeout    = 3 * time.Second
	semaphoreSlots = uint64(1)
)

func init() {
	locksmith2Cmd.AddCommand(cmdDaemon)
}

func runDaemon(cmd *cobra.Command, cmdArgs []string) error {
	listenAddr := fmt.Sprintf("%s:%d", address, port)
	logrus.WithFields(logrus.Fields{
		"address": address,
		"port":    port,
	}).Info("starting daemon")

	serverConf := daemon.ServerConfig{
		EtcdURLs:       etcdURLs,
		LockTimeout:    lockTimeout,
		SemaphoreSlots: semaphoreSlots,
	}

	http.Handle(daemon.PreRebootEndpoint, serverConf.PreReboot())
	http.Handle(daemon.SteadyStateEndpoint, serverConf.SteadyState())
	return http.ListenAndServe(listenAddr, nil)
}
