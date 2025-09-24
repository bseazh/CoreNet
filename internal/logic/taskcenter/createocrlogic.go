package taskcenter

import (
	"context"

	"corenet/internal/svc"
	"corenet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOCRLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateOCRLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOCRLogic {
	return &CreateOCRLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateOCRLogic) CreateOCR(req *types.CreateOCRReq) (resp *types.CreateOCRResp, err error) {
	// todo: add your logic here and delete this line

	return
}
