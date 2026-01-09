package fkit

import stdlog "log"

type Logger interface {
	Debug(args ...any)
	Debugf(format string, args ...any)
	Info(args ...any)
	Infof(format string, args ...any)
	Warn(args ...any)
	Warnf(format string, args ...any)
	Error(args ...any)
	Errorf(format string, args ...any)
	Fatal(args ...any)
	Fatalf(format string, args ...any)
}

type StdLogger struct {
}

func (s *StdLogger) Debug(args ...any) {
	stdlog.Default().Println(args...)
}

func (s *StdLogger) Debugf(format string, args ...any) {
	stdlog.Default().Printf(format+"\n", args...)
}

func (s *StdLogger) Info(args ...any) {
	stdlog.Default().Println(args...)
}

func (s *StdLogger) Infof(format string, args ...any) {
	stdlog.Default().Printf(format+"\n", args...)
}

func (s *StdLogger) Warn(args ...any) {
	stdlog.Default().Println(args...)
}

func (s *StdLogger) Warnf(format string, args ...any) {
	stdlog.Default().Printf(format+"\n", args...)
}

func (s *StdLogger) Error(args ...any) {
	stdlog.Default().Println(args...)
}

func (s *StdLogger) Errorf(format string, args ...any) {
	stdlog.Default().Printf(format+"\n", args...)
}

func (s *StdLogger) Fatal(args ...any) {
	stdlog.Default().Println(args...)
}

func (s *StdLogger) Fatalf(format string, args ...any) {
	stdlog.Default().Printf(format+"\n", args...)
}
