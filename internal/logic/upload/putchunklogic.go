package upload

import (
	"context"

	"corenet/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type PutChunkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPutChunkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PutChunkLogic {
	return &PutChunkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PutChunkLogic) PutChunk() error {
	// todo: add your logic here and delete this line

	return nil
}
