package application

import (
	"context"
	"errors"
	"time"
)

// DeleteOldStates to clear the database
func (a *Application) DeleteOldStates(ctx context.Context) error {
	// delete old states
	threshold := time.Now().AddDate(0, 0, -4)
	err := a.repo.DeleteOldStates(ctx, threshold)
	if err != nil {
		a.log.Error(err)
		return errors.New("internal server error")
	}

	a.log.Trace("old states deleted")
	return nil
}
