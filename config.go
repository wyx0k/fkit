package fkit

import (
	"github.com/spf13/viper"
	"github.com/wyx0k/fkit/config"
)

func (a *app[T]) ParseConfig(path string, t *T, opts ...viper.DecoderConfigOption) {
	a.Init(func(ctx *FkitContext[T]) (CloseHook[T], error) {
		ctx.conf = t
		err := config.InitConfig(path, t, opts...)
		if err != nil {
			return nil, err
		}

		return nil, err
	})
}

func (a *app[T]) ParseConfigWithDecrypt(path string, t *T, decryptFunc func(t *T) error, opts ...viper.DecoderConfigOption) {
	a.Init(func(ctx *FkitContext[T]) (CloseHook[T], error) {
		ctx.conf = t
		err := config.InitConfig(path, t, opts...)
		if err != nil {
			return nil, err
		}
		err = decryptFunc(t)
		return nil, err
	})
}

func (a *app[T]) WithConfig(t *T) {
	a.Init(func(ctx *FkitContext[T]) (CloseHook[T], error) {
		ctx.conf = t
		return nil, nil
	})
}
