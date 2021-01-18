// Code generated by goctl. DO NOT EDIT!
// Source: check.proto

//go:generate mockgen -destination ./check_mock.go -package checkclient -source $GOFILE

package checkclient

import (
	"context"

	"github.com/notes/go_zero/rpc/check/check"

	"github.com/tal-tech/go-zero/zrpc"
)

type (
	Request  = check.Request
	Response = check.Response

	Check interface {
		Ping(ctx context.Context, in *Request) (*Response, error)
	}

	defaultCheck struct {
		cli zrpc.Client
	}
)

func NewCheck(cli zrpc.Client) Check {
	return &defaultCheck{
		cli: cli,
	}
}

func (m *defaultCheck) Ping(ctx context.Context, in *Request) (*Response, error) {
	client := check.NewCheckClient(m.cli.Conn())
	return client.Ping(ctx, in)
}
