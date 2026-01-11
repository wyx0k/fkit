package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	app "github.com/wyx0k/fkit"

	log "github.com/op/go-logging"
	"gopkg.in/natefinch/lumberjack.v2"
)

var format = log.MustStringFormatter(
	`%{color}%{time:2006-01-02 15:04:05.000} %{level:.5s} %{shortfile}:%{color:reset} %{message}`,
)

type LogConf struct {
	Path       string `json:"path" yaml:"path"`
	MaxSize    int    `json:"max_size" yaml:"max_size"`
	MaxAge     int    `json:"max_age" yaml:"max_age"`
	MaxBackups int    `json:"max_backups" yaml:"max_backups"`
	Compress   bool   `json:"compress" yaml:"compress"`
}

type LogConfHolder interface {
	LogConf() LogConf
}

func SimpleLog[T LogConfHolder](ctx *app.FkitContext[T]) (app.CloseHook[T], error) {
	conf := (*ctx.Conf()).LogConf()
	if conf.Path == "" {
		conf.Path = "./"
	}
	if conf.MaxSize == 0 {
		conf.MaxSize = 1024 * 1024 * 100
	}
	if conf.MaxAge == 0 {
		conf.MaxAge = 30
	}
	sl := &SimpleLogger{Conf: conf, AppName: ctx.AppName()}
	path := sl.Init()
	ctx.SetAppLog(sl)
	_logger = sl
	Infof("日志初始化完毕,日志位置:%s", path)
	return nil, nil
}

type SimpleLogger struct {
	AppName    string
	Conf       LogConf
	logger     *log.Logger
	lumberjack *lumberjack.Logger
}

func (s *SimpleLogger) Init() string {
	path := fmt.Sprintf("%s/%s.log", s.Conf.Path, s.AppName)
	p, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	s.lumberjack = &lumberjack.Logger{
		Filename:   p,
		MaxSize:    s.Conf.MaxSize,
		MaxAge:     s.Conf.MaxAge,
		MaxBackups: s.Conf.MaxBackups,
		Compress:   s.Conf.Compress,
		LocalTime:  true,
	}
	out := io.MultiWriter(s.lumberjack, os.Stdout)
	backend := log.NewLogBackend(out, "", 0)
	backendFormatter := log.NewBackendFormatter(backend, format)
	log.SetBackend(backendFormatter)
	s.logger = log.MustGetLogger(s.AppName)
	s.logger.ExtraCalldepth = 2
	_logger = s
	return p
}

func (s *SimpleLogger) Debug(args ...any) {
	s.logger.Debug(args...)
}

func (s *SimpleLogger) Debugf(format string, args ...any) {
	s.logger.Debugf(format, args...)
}

func (s *SimpleLogger) Info(args ...any) {
	s.logger.Info(args...)
}

func (s *SimpleLogger) Infof(format string, args ...any) {
	s.logger.Infof(format, args...)
}

func (s *SimpleLogger) Warn(args ...any) {
	s.logger.Warning(args...)
}

func (s *SimpleLogger) Warnf(format string, args ...any) {
	s.logger.Warningf(format, args...)
}

func (s *SimpleLogger) Error(args ...any) {
	s.logger.Error(args...)
}

func (s *SimpleLogger) Errorf(format string, args ...any) {
	s.logger.Errorf(format, args...)
}

func (s *SimpleLogger) Fatal(args ...any) {
	s.logger.Fatal(args...)
}

func (s *SimpleLogger) Fatalf(format string, args ...any) {
	s.logger.Fatalf(format, args...)
}
