package application

import (
	"context"
	"errors"
	"time"
)

// DeleteOldTrainAlarms to clear the database
func (a *Application) DeleteOldTrainAlarms(ctx context.Context) error {
	// delete old trainalarms
	threshold := time.Now().AddDate(0, 0, -2)
	a.log.Infof("delete old trains before %v", threshold)
	err := a.repo.DeleteOldTrainAlarms(ctx, threshold)
	if err != nil {
		a.log.Error(err)
		return errors.New("internal server error")
	}

	a.log.Trace("old train alarms deleted")
	return nil
}
