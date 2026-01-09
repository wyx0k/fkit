package fkit

import (
	"context"
	"os"
	"os/signal"
	"slices"
	"syscall"
)

type InitHook[T any] func(ctx *FkitContext[T]) (CloseHook[T], error)
type CloseHook[T any] func(ctx *FkitContext[T], code int) error
type PaincHook[T any] func(ctx *FkitContext[T], v any) int
type ReloadHook[T any] func(ctx *FkitContext[T]) error

type app[T any] struct {
	name        string
	version     string
	ctx         *FkitContext[T]
	initHooks   []InitHook[T]
	closeHooks  []CloseHook[T]
	panicHook   PaincHook[T]
	reloadHook  ReloadHook[T]
	appSigCh    chan os.Signal
	appSelfExit chan string
}

func NewApp[T any](name string, verion string, ctx ...context.Context) *app[T] {
	app := &app[T]{name: name, version: verion}
	originalCtx := context.TODO()
	if len(ctx) != 0 {
		originalCtx = ctx[0]
	}
	fCtx := &FkitContext[T]{app: app, ctx: originalCtx, values: map[string]any{}, appLog: &StdLogger{}}
	app.ctx = fCtx
	app.appSigCh = make(chan os.Signal, 1)
	app.appSelfExit = make(chan string, 1)
	signal.Notify(app.appSigCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	return app
}
func (a *app[T]) SetBanner(banner string) {
	a.ctx.banner = banner
}
func (a *app[T]) Init(fns ...InitHook[T]) {
	a.initHooks = append(a.initHooks, fns...)
}

func (a *app[T]) InitPanic(fn PaincHook[T]) {
	a.panicHook = fn
}

func (a *app[T]) Reload(fn ReloadHook[T]) {
	a.reloadHook = fn
}

func (a *app[T]) Start() {
	a.start(nil)
}

func (a *app[T]) StartAsJob(job func(ctx *FkitContext[T]) error) {
	a.start(job)
}

func (a *app[T]) start(job func(ctx *FkitContext[T]) error) {
	code := 0
	defer func() {
		a.ctx.changePhase(AppFinishing)
		if v := recover(); v != nil {
			if a.panicHook != nil {
				os.Exit(a.panicHook(a.ctx, v))
			} else {
				panic(v)
			}
		} else {
			os.Exit(code)
		}
	}()
	a.ctx.changePhase(AppInitializing)
	for _, hook := range a.initHooks {
		closeFn, err := hook(a.ctx)
		if err != nil {
			a.ctx.appLog.Error("init failed: ", err)
			os.Exit(1)
		}
		if closeFn != nil {
			a.closeHooks = append(a.closeHooks, closeFn)
		}
	}
	slices.Reverse(a.closeHooks)
	a.ctx.changePhase(AppRunning)
	defer func() {
		a.ctx.changePhase(AppClosing)
		for _, hook := range a.closeHooks {
			err := hook(a.ctx, code)
			if err != nil {
				a.ctx.appLog.Error("close failed: ", err)
			}
		}
	}()
	if job != nil {
		go func() {
			err := job(a.ctx)
			if err != nil {
				a.ctx.Exit(err.Error())
			} else {
				a.ctx.Exit("job done")
			}
		}()
	}
	code = a.waitForDone(a.reloadHook)

}

func (a *app[T]) waitForDone(reloadFn ReloadHook[T]) int {
	for {
		select {
		case sig := <-a.appSigCh:
			switch sig {
			case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				return 0
			case syscall.SIGHUP:
				// do nothing
				if reloadFn != nil {
					reloadFn(a.ctx)
				}
			default:
				return 1
			}
		case msg := <-a.appSelfExit:
			a.ctx.appLog.Info("app自行退出,msg：", msg)
			return 0
		}
	}
}

func (a *app[T]) exit(msg string) {
	a.appSelfExit <- msg
}
