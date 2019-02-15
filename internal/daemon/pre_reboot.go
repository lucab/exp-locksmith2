package daemon

import (
	"context"
	"net/http"
	"time"

	"github.com/lucab/exp-locksmith2/internal/lock"
	"github.com/sirupsen/logrus"
)

// PreReboot is the handler for the `/v1/pre-reboot` endpoint.
func PreReboot(w http.ResponseWriter, req *http.Request) {
	logrus.Debug("got pre-reboot request")

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

	err = lockManager.RecursiveLock(ctx, nodeIdentity.Group)
	if err != nil {
		logrus.Errorln(err)
		http.Error(w, err.Error(), 500)
		return
	}

	logrus.Debug("green-flag to pre-reboot request")
}
