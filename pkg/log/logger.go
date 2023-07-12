package log

import "go.uber.org/zap"

func NewLogger() (*zap.SugaredLogger, error) {
	l, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	return l.Sugar(), nil
}

func GlobalLogger() *zap.SugaredLogger {
	return zap.L().Sugar()
}
