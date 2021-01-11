package ratelimit

import "context"

// Op ...
type Op int

const (
	//Success ...
	Success Op = iota
	//Ignore ...
	Ignore
	//Drop ...
	Drop
)

type allowOptions struct{}

//AllowOption ...
type AllowOption interface {
	Apply(*allowOptions)
}

//DoneInfo ...
type DoneInfo struct {
	Err error
	Op  Op
}

//DefaultAllowOpts ..
func DefaultAllowOpts() allowOptions {
	return allowOptions{}
}

// Limiter limit interface.
type Limiter interface {
	Allow(ctx context.Context, opts ...AllowOption) (func(info DoneInfo), error)
}
