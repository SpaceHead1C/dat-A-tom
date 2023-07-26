package routines

import "time"

func NewSendChangedDataRoutine(c SendChangedDataConfig) func() error {
	return func() error {
		c.Logger.Infof("routine started at %s", time.Now())
		return nil
	}
}
