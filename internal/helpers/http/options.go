package http

import "time"

var (
	Opts *Options
)

type Option func(*Options)

type Options struct {
	Tls
	Timeout time.Duration
	Retry
}

type Tls struct {
	InsecureSkipVerify bool
}

type Retry struct {
	Count       int
	WaitTime    time.Duration
	MaxWaitTime time.Duration
}

func TlsInsecureSkipVerify(skip bool) Option {
	return func(o *Options) {
		o.Tls.InsecureSkipVerify = skip
	}
}

func Timeout(t time.Duration) Option {
	return func(o *Options) {
		o.Timeout = t
	}
}

func RetryCount(c int) Option {
	return func(o *Options) {
		o.Retry.Count = c
	}
}

func RetryWaitTime(t time.Duration) Option {
	return func(o *Options) {
		o.Retry.WaitTime = t
	}
}

func RetryMaxWaitTime(t time.Duration) Option {
	return func(o *Options) {
		o.Retry.MaxWaitTime = t
	}
}
