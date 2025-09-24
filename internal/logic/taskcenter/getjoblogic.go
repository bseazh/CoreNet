package taskcenter

import (
	"context"

	"corenet/internal/svc"
	"corenet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetJobLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetJobLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetJobLogic {
	return &GetJobLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetJobLogic) GetJob() (resp *types.JobResp, err error) {
	// todo: add your logic here and delete this line

	return
}
