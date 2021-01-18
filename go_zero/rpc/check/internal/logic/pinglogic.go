package logic

import (
	"context"

	"github.com/notes/go_zero/rpc/check/check"
	"github.com/notes/go_zero/rpc/check/internal/svc"

	"github.com/tal-tech/go-zero/core/logx"
)

type PingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PingLogic {
	return &PingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PingLogic) Ping(in *check.Request) (*check.Response, error) {
	// todo: add your logic here and delete this line

	return &check.Response{}, nil
}
