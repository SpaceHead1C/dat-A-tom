package routines

func NewSendChangedDataRoutine(c SendChangedDataConfig) func() error {
	return func() error {
		err := sendChangedData(c)
		if err != nil {
			c.Logger.Errorln(err.Error())
		}
		return err
	}
}
