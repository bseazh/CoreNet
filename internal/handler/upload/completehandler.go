package upload

import (
	"net/http"

	"corenet/internal/logic/upload"
	"corenet/internal/svc"
	"corenet/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func CompleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CompleteReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := upload.NewCompleteLogic(r.Context(), svcCtx)
		resp, err := l.Complete(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
