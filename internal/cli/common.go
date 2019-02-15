package cli

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	// locksmith2Cmd is the top-level cobra command for `torcx`
	locksmith2Cmd = &cobra.Command{
		Use:           "locksmith2",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
)

// Init initializes the CLI environment for locksmith2
func Init() (*cobra.Command, error) {
	logrus.SetLevel(logrus.DebugLevel)
	return locksmith2Cmd, nil
}
