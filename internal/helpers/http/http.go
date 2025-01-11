package http

import (
	"crypto/tls"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

var (
	once   sync.Once
	helper *Helper
)

type Client interface {
	R() *resty.Request
}

type Helper struct {
	Client
	Options
}

func initOptions(opts []Option) *Options {
	options := genDefaultOptions()
	for _, o := range opts {
		o(options)
	}

	return options
}

func genDefaultOptions() *Options {
	return &Options{
		Tls: Tls{
			InsecureSkipVerify: true,
		},
		Timeout: 10 * time.Second,
		Retry: Retry{
			Count:       3,
			WaitTime:    2 * time.Second,
			MaxWaitTime: 5 * time.Second,
		},
	}
}

func NewHelper(opts ...Option) (*Helper, error) {
	r := resty.New()
	initedOpts := initOptions(opts)

	r.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: initedOpts.Tls.InsecureSkipVerify})
	r.SetTimeout(initedOpts.Timeout)
	r.SetRetryCount(initedOpts.Retry.Count)
	r.SetRetryWaitTime(initedOpts.Retry.WaitTime)
	r.SetRetryMaxWaitTime(initedOpts.Retry.MaxWaitTime)

	return &Helper{
		Client:  r,
		Options: *initedOpts,
	}, nil
}

func NewGlobalHelper(opts ...Option) error {
	var h *Helper
	var err error

	once.Do(func() {
		h, err = NewHelper(opts...)
		if err != nil {
			return
		}

		helper = h
	})
	if err != nil {
		return err
	}

	return nil
}

func GetGlobalHelper() *Helper {
	return helper
}
