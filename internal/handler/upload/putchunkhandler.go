package upload

import (
	"net/http"

	"corenet/internal/logic/upload"
	"corenet/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func PutChunkHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := upload.NewPutChunkLogic(r.Context(), svcCtx)
		err := l.PutChunk()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
