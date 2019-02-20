package daemon

import (
	"context"
	"net/http"

	"github.com/lucab/exp-locksmith2/internal/lock"
	"github.com/sirupsen/logrus"
)

const (
	// PreRebootEndpoint is the endpoint for requesting a semaphore lock.
	PreRebootEndpoint = "/v1/pre-reboot"
)

// PreReboot is the handler for the `/v1/pre-reboot` endpoint.
func (sc *ServerConfig) PreReboot() http.Handler {
	handler := func(w http.ResponseWriter, req *http.Request) {
		if sc == nil {
			http.Error(w, errNilServerConfig.Error(), 500)
			return
		}

		logrus.Debug("got pre-reboot request")

		nodeIdentity, err := validateIdentity(req)
		if err != nil {
			logrus.Errorln("failed to validate client identity: ", err)
			http.Error(w, err.Error(), 400)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), sc.LockTimeout)
		defer cancel()
		lockManager, err := lock.NewManager(ctx, sc.EtcdURLs, nodeIdentity.Group, sc.SemaphoreSlots)
		if err != nil {
			logrus.Errorln("failed to initialize semaphore manager: ", err)
			http.Error(w, err.Error(), 500)
			return
		}
		defer lockManager.Close()

		err = lockManager.RecursiveLock(ctx, nodeIdentity.UUID)
		if err != nil {
			logrus.Errorln(err)
			http.Error(w, err.Error(), 500)
			return
		}

		logrus.WithFields(logrus.Fields{
			"group": nodeIdentity.Group,
			"UUID":  nodeIdentity.UUID,
		}).Debug("green-flag to pre-reboot request")
	}

	return http.HandlerFunc(handler)
}
