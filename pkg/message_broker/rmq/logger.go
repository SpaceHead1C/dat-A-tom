package rmq

import "go.uber.org/zap"

type logger struct {
	l *zap.SugaredLogger
}

func (l logger) Fatalf(template string, args ...any) {
	l.l.Fatalf(template, args)
}

func (l logger) Errorf(template string, args ...any) {
	l.l.Errorf(template, args)
}

func (l logger) Warnf(template string, args ...any) {
	l.l.Warnf(template, args)
}

func (l logger) Infof(template string, args ...any) {
	l.l.Infof(template, args)
}

func (l logger) Debugf(template string, args ...any) {
	l.l.Debugf(template, args)
}

func (l logger) Tracef(template string, args ...any) {
	l.l.Infof(template, args)
}
