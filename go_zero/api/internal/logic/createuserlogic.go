package logic

import (
	"context"

	"github.com/notes/go_zero/api/internal/svc"
	"github.com/notes/go_zero/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type CreateUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) CreateUserLogic {
	return CreateUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateUserLogic) CreateUser(req types.Request) error {
	// todo: add your logic here and delete this line

	return nil
}
