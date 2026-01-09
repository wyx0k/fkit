package fkit

import "github.com/wyx0k/fkit/config"

func (a *app[T]) ParseConfig(path string, t *T, decryptFunc ...func(t *T) error) {
	a.Init(func(ctx *FkitContext[T]) (CloseHook[T], error) {
		ctx.conf = t
		err := config.InitConfig(path, t)
		if err != nil {
			return nil, err
		}
		if len(decryptFunc) > 0 {
			err = decryptFunc[0](t)
		}
		return nil, err
	})
}
