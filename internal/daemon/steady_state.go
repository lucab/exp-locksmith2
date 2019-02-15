package daemon

import (
	"context"
	"net/http"
	"time"

	"github.com/lucab/exp-locksmith2/internal/lock"
	"github.com/sirupsen/logrus"
)

// SteadyState is the handler for the `/v1/steady-state` endpoint.
func SteadyState(w http.ResponseWriter, req *http.Request) {
	logrus.Debug("got steady-state report")

	nodeIdentity, err := validateIdentity(req)
	if err != nil {
		logrus.Errorln(err)
		http.Error(w, err.Error(), 400)
		return
	}

	ctx := context.TODO()
	lockManager, err := lock.NewManager(ctx, 5*time.Second, nodeIdentity.UUID)
	if err != nil {
		logrus.Errorln(err)
		http.Error(w, err.Error(), 500)
		return
	}
	defer lockManager.Close()

	err = lockManager.UnlockIfHeld(ctx, nodeIdentity.Group)
	if err != nil {
		logrus.Errorln(err)
		http.Error(w, err.Error(), 500)
		return
	}

	logrus.Debug("steady-state confirmed")
}
