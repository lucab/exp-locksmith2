package cli

import (
	"fmt"
	"net/http"
	"time"

	"github.com/lucab/exp-locksmith2/internal/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	cmdServe = &cobra.Command{
		Use:  "serve",
		RunE: runServe,
	}
	address        = "0.0.0.0"
	port           = 9999
	etcdURLs       = []string{"http://127.0.0.1:2379"}
	lockTimeout    = 3 * time.Second
	semaphoreSlots = uint64(1)
)

func init() {
	locksmith2Cmd.AddCommand(cmdServe)
}

func runServe(cmd *cobra.Command, cmdArgs []string) error {
	listenAddr := fmt.Sprintf("%s:%d", address, port)
	logrus.WithFields(logrus.Fields{
		"address": address,
		"port":    port,
	}).Info("starting service")

	config := server.ServerConfig{
		EtcdURLs:       etcdURLs,
		LockTimeout:    lockTimeout,
		SemaphoreSlots: semaphoreSlots,
	}

	http.Handle(server.PreRebootEndpoint, config.PreReboot())
	http.Handle(server.SteadyStateEndpoint, config.SteadyState())
	return http.ListenAndServe(listenAddr, nil)
}
