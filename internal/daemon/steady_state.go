package daemon

import (
	"context"
	"net/http"

	"github.com/lucab/exp-locksmith2/internal/lock"
	"github.com/sirupsen/logrus"
)

const (
	// SteadyStateEndpoint is the endpoint for releasing a semaphore lock.
	SteadyStateEndpoint = "/v1/steady-state"
)

// SteadyState is the handler for the `/v1/steady-state` endpoint.
func (sc *ServerConfig) SteadyState() http.Handler {
	handler := func(w http.ResponseWriter, req *http.Request) {
		if sc == nil {
			http.Error(w, errNilServerConfig.Error(), 500)
			return
		}

		logrus.Debug("got steady-state report")
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

		err = lockManager.UnlockIfHeld(ctx, nodeIdentity.UUID)
		if err != nil {
			logrus.Errorln("failed to release any semaphore lock: ", err)
			http.Error(w, err.Error(), 500)
			return
		}

		logrus.WithFields(logrus.Fields{
			"group": nodeIdentity.Group,
			"UUID":  nodeIdentity.UUID,
		}).Debug("steady-state confirmed")
	}

	return http.HandlerFunc(handler)
}
