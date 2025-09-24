package preview

import (
	"context"

	"corenet/internal/svc"
	"corenet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPreviewLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPreviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPreviewLogic {
	return &GetPreviewLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPreviewLogic) GetPreview() (resp *types.PreviewResp, err error) {
	// todo: add your logic here and delete this line

	return
}
