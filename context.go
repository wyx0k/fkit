package fkit

import (
	"context"
)

const defaultTitle = `
  ______   _  __  _____   _______ 
 |  ____| | |/ / |_   _| |__   __|
 | |__    | ' /    | |      | |   
 |  __|   |  <     | |      | |   
 | |      | . \   _| |_     | |   
 |_|      |_|\_\ |_____|    |_|   
                                                                 
`

type AppPhase string

const (
	AppInitializing AppPhase = "initializing"
	AppRunning      AppPhase = "running"
	AppClosing      AppPhase = "closing"
	AppFinishing    AppPhase = "finishing"
)

type GoroutinePaincHook[T any] func(ctx *FkitContext[T], v any) bool

type FkitContext[T any] struct {
	app    *app[T]
	banner string
	conf   *T
	ctx    context.Context
	phase  AppPhase
	values map[string]any
	appLog Logger
}

func (fc *FkitContext[T]) changePhase(phase AppPhase) {
	fc.appLog.Info("app_phase-> ", phase)
	fc.phase = phase
	if phase == AppRunning {
		if fc.banner == "" {
			fc.appLog.Infof("%s\nApp Version: %s", defaultTitle, fc.app.version)
		} else {
			fc.appLog.Infof("%s\nApp Version: %s", fc.banner, fc.app.version)
		}
	}
}
func (fc *FkitContext[T]) AppName() string {
	return fc.app.name
}

func (fc *FkitContext[T]) Exit(msg string) {
	if fc.phase != AppRunning {
		return
	}
	fc.app.exit(msg)
}

func (fc *FkitContext[T]) Ctx() context.Context {
	return fc.ctx
}
func (fc *FkitContext[T]) Conf() *T {
	return fc.conf
}

func (fc *FkitContext[T]) SetAppLog(logger Logger) {
	fc.appLog = logger
}

func (fc *FkitContext[T]) Set(name string, value any) {
	fc.values[name] = value
}
func (fc *FkitContext[T]) Get(name string) any {
	return fc.values[name]
}
func (fc *FkitContext[T]) SafeGo(fn func(), name ...string) {
	go func() {
		defer func() {
			if v := recover(); v != nil {
				if len(name) > 0 {
					fc.appLog.Errorf("协程[%s]因panic退出:%v", name[0], v)
				} else {
					fc.appLog.Errorf("协程因panic退出:%v", v)
				}
			}
		}()
		fn()
	}()
}

func (fc *FkitContext[T]) GoWithPanicHook(fn func(), hook GoroutinePaincHook[T]) {
	go func() {
		defer func() {
			if v := recover(); v != nil {
				if hook != nil {
					restart := hook(fc, v)
					if restart {
						fc.GoWithPanicHook(fn, hook)
					}
				}
			}
		}()
		fn()
	}()
}
